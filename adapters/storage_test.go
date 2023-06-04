package adapters

import (
	"github.com/google/uuid"
	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlStorageTestSuite struct {
	suite.Suite
	testTeam   *domain.Team
	testMember *domain.TeamMember
	db         *gorm.DB
}

func (suite *SqlStorageTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Nil(err)

	err = db.AutoMigrate(&domain.Team{}, &domain.TeamMember{})
	suite.Nil(err)

	suite.db = db
	suite.testMember = &domain.TeamMember{
		UUID:    uuid.NewString(),
		Email:   "fake@example.com",
		IsAdmin: true,
		IsOwner: true,
	}
	suite.testTeam = &domain.Team{
		UUID:        uuid.NewString(),
		Name:        "Test team",
		Description: "Test create team",
		Members:     []domain.TeamMember{*suite.testMember},
	}
}

func (suite *SqlStorageTestSuite) TestDBConnection() {
	result := suite.db.Create(suite.testTeam)
	suite.Nil(result.Error)
	suite.Equal(int64(1), result.RowsAffected)
}

func (suite *SqlStorageTestSuite) TestAddWithSqlRepository() {
	repo := GormSqlRepository{suite.db}

	err := repo.Add(suite.testTeam)
	suite.Nil(err)

	fetchTeam := new(domain.Team)
	suite.db.Take(fetchTeam)
	suite.Equal(suite.testTeam.Name, fetchTeam.Name)

	fetchMember := new(domain.TeamMember)
	suite.db.Take(fetchMember)
	suite.Equal(suite.testMember.UUID, fetchMember.UUID)
	suite.Equal(fetchMember.TeamUUID, fetchTeam.UUID)
	suite.True(fetchMember.IsAdmin)
	suite.True(fetchMember.IsOwner)
}

func (suite *SqlStorageTestSuite) TestGetWithSqlRepository() {
	result := suite.db.Create(suite.testTeam)
	suite.Nil(result.Error)

	repo := GormSqlRepository{suite.db}

	fetchTeam, err := repo.Get(suite.testTeam.UUID)
	suite.Nil(err)
	suite.Equal(suite.testTeam.UUID, fetchTeam.UUID)
	suite.Equal(1, len(fetchTeam.Members))
	suite.Equal(suite.testMember.UUID, fetchTeam.Members[0].UUID)
}

func (suite *SqlStorageTestSuite) TestCreateWithSqlUOW() {
	uow := &GormSqlUnitOfWork{suite.db, nil, nil}
	command := &domain.CreateTeamCommand{
		Name:        "Test team",
		Description: "Test create team",
		OwnerEmail:  "fake@example.com",
	}

	teamUUID, err := service_layer.CreateTeam(command, uow)
	suite.Nil(err)

	fetchTeam := new(domain.Team)
	suite.db.First(fetchTeam, "UUID = ?", teamUUID)
	suite.Equal(command.Name, fetchTeam.Name)
	suite.Equal(len(fetchTeam.Members), 1)
	suite.Equal(fetchTeam.Members[0].Email, command.OwnerEmail)
}

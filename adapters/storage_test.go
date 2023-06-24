package adapters

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlStorageTestSuite struct {
	suite.Suite
	testTeam   *domain.Team
	testMember *domain.TeamMember
	db         *gorm.DB
}

func (suite *SqlStorageTestSuite) SetupSuite() {
	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Nil(err)

	suite.db = db
}

func (suite *SqlStorageTestSuite) SetupTest() {
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

func (suite *SqlStorageTestSuite) TearDownTest() {
	suite.db.Exec("TRUNCATE teams CASCADE")
}

func (suite *SqlStorageTestSuite) TestDbConnection() {
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
	suite.Equal(suite.testTeam.UUID, fetchTeam.UUID)

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
	uow := &GormSqlUnitOfWork{db: suite.db}
	command := &domain.CreateTeamCommand{
		Name:        "Test team",
		Description: "Test create team",
		OwnerEmail:  "fake@example.com",
	}

	teamUUID, err := service_layer.CreateTeam(command, uow)
	suite.Nil(err)

	fetchTeam := new(domain.Team)
	result := suite.db.Model(fetchTeam).Preload("Members").First(fetchTeam, "uuid = ?", teamUUID)
	suite.Nil(result.Error)
	suite.Equal(command.Name, fetchTeam.Name)
	suite.Equal(1, len(fetchTeam.Members))
	suite.Equal(fetchTeam.Members[0].Email, command.OwnerEmail)
}

func (suite *SqlStorageTestSuite) TestUOWRollback() {
	uow := &GormSqlUnitOfWork{db: suite.db}
	uow.Begin()
	uow.Teams().Add(suite.testTeam)
	uow.Rollback()

	var teams []domain.Team
	result := suite.db.Find(&teams)
	suite.Nil(result.Error)
	suite.Equal(int64(0), result.RowsAffected)
}

func (suite *SqlStorageTestSuite) TestUOWCommitAfterRollbackFails() {
	uow := &GormSqlUnitOfWork{db: suite.db}
	uow.Begin()
	uow.Teams().Add(suite.testTeam)
	uow.Rollback()

	uow.Teams().Add(suite.testTeam)
	err := uow.Commit()
	suite.Error(err)

	var teams []domain.Team
	result := suite.db.Find(&teams)
	suite.Nil(result.Error)
	suite.Equal(int64(0), result.RowsAffected)
}

func (suite *SqlStorageTestSuite) TestUOWConcurrentWrites() {
	var wg sync.WaitGroup
	names := [...]string{"Test1", "Test2"}

	fc := func(name string) {
		defer wg.Done()
		uow := &GormSqlUnitOfWork{db: suite.db}
		uow.Begin()
		uow.Teams().Add(&domain.Team{UUID: uuid.NewString(), Name: name})
		uow.Commit()
	}

	for _, teamName := range names {
		wg.Add(1)
		go fc(teamName)
	}
	wg.Wait()

	var teams []domain.Team
	result := suite.db.Find(&teams)
	suite.Nil(result.Error)
	suite.Equal(int64(2), result.RowsAffected)
}

func TestSqlStorageTestSuite(t *testing.T) {
	suite.Run(t, new(SqlStorageTestSuite))
}

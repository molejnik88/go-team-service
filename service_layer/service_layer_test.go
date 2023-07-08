package service_layer

import (
	"fmt"
	"testing"

	"github.com/molejnik88/go-team-service/domain"
	"github.com/stretchr/testify/assert"
)

type FakeRepository struct {
	teams map[string]domain.Team
}

func (r *FakeRepository) Add(team *domain.Team) error {
	r.teams[team.UUID] = *team
	return nil
}

func (r *FakeRepository) Get(uuid string) (*domain.Team, error) {
	if team, ok := r.teams[uuid]; ok {
		return &team, nil
	}

	return nil, fmt.Errorf("Team with uuid: %s does not exist", uuid)
}

type FakeUOW struct {
	repository *FakeRepository
	commited   bool
	rollbacked bool
}

func (uow *FakeUOW) Begin() error {
	return nil
}

func (uow *FakeUOW) Teams() Repository {
	return uow.repository
}

func (uow *FakeUOW) Commit() error {
	uow.commited = true
	return nil
}

func (uow *FakeUOW) Rollback() {
	uow.rollbacked = true
}

func TestDummyTeamName(t *testing.T) {
	expectedTeamName := "Test Team"
	testTeam := domain.Team{
		Name: "Test Team",
	}

	assert.Equal(t, expectedTeamName, testTeam.Name, fmt.Sprintf("Incorrect team name; expected \"%s\"", expectedTeamName))
}

func TestCreateTeam(t *testing.T) {
	createCommand := &domain.CreateTeamCommand{
		Name:        "Test Team",
		Description: "Test Description",
	}
	uow := &FakeUOW{
		repository: &FakeRepository{
			teams: make(map[string]domain.Team),
		},
		commited:   false,
		rollbacked: false,
	}

	teamUUID, err := CreateTeam(createCommand, uow)
	assert.NotEmpty(t, teamUUID)
	assert.Nil(t, err)
	assert.True(t, uow.commited)
	assert.True(t, uow.rollbacked)

	newTeam, err := uow.Teams().Get(teamUUID)
	assert.Nil(t, err)
	assert.Equal(t, newTeam.Name, "Test Team")
}

func TestCreateTeamSetsOwnerCorrectly(t *testing.T) {
	createCommand := &domain.CreateTeamCommand{
		Name:        "Test Team",
		Description: "Test Description",
		OwnerEmail:  "fake@example.com",
	}
	uow := &FakeUOW{
		repository: &FakeRepository{
			teams: make(map[string]domain.Team),
		},
		commited:   false,
		rollbacked: false,
	}

	teamUUID, err := CreateTeam(createCommand, uow)
	assert.Nil(t, err)

	newTeam, err := uow.Teams().Get(teamUUID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(newTeam.Members))
	newTeamMember := newTeam.Members[0]
	assert.True(t, newTeamMember.IsAdmin)
	assert.True(t, newTeamMember.IsOwner)
}

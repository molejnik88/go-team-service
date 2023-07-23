package service_layer

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/molejnik88/go-team-service/domain"
)

func CreateTeam(command *domain.CreateTeamCommand, uow UnitOfWork) (string, error) {
	uow.Begin()
	defer uow.Rollback()

	newTeam := &domain.Team{
		UUID:        uuid.NewString(),
		Name:        command.Name,
		Description: command.Description,
		Members: []domain.TeamMember{
			{UUID: uuid.NewString(), Email: command.OwnerEmail, IsAdmin: true, IsOwner: true},
		},
	}

	err := uow.Teams().Add(newTeam)
	if err != nil {
		return "", err
	}

	err = uow.Commit()
	if err != nil {
		return "", err
	}

	return newTeam.UUID, nil
}

func AddTeamMember(command *domain.AddTeamMember, uow UnitOfWork) error {
	uow.Begin()
	defer uow.Rollback()

	team, err := uow.Teams().Get(command.TeamUUID)
	if err != nil {
		return fmt.Errorf("could not fetch team with uuid: %s", command.TeamUUID)
	}

	m := domain.TeamMember{
		UUID:    uuid.NewString(),
		Email:   command.Email,
		IsAdmin: command.IsAdmin,
	}

	err = team.AddMember(m)
	if err != nil {
		return err
	}

	uow.Teams().Update(team)
	return uow.Commit()
}

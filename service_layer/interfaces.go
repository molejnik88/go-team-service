package service_layer

import "github.com/molejnik88/go-team-service/domain"

type Repository interface {
	Get(uuid string) (*domain.Team, error)
	Add(team *domain.Team) error
}

type UnitOfWork interface {
	Teams() Repository
	Begin() error
	Commit() error
	Rollback()
}

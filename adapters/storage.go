package adapters

import (
	"fmt"

	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
	"gorm.io/gorm"
)

// TODO: Remove when a db implementation is ready
type InMemoryRepository struct {
	Teams map[string]domain.Team
}

func (r *InMemoryRepository) Add(team *domain.Team) error {
	r.Teams[team.UUID] = *team
	return nil
}

func (r *InMemoryRepository) Get(uuid string) (*domain.Team, error) {
	if team, ok := r.Teams[uuid]; ok {
		return &team, nil
	}

	return nil, fmt.Errorf("team with uuid: %s does not exist", uuid)
}

// TODO: Remove when a db implementation is ready
type InMemoryUOW struct {
	Repository *InMemoryRepository
	commited   bool
	rollbacked bool
}

func (uow *InMemoryUOW) Begin() error {
	return nil
}

func (uow *InMemoryUOW) Teams() service_layer.Repository {
	return uow.Repository
}

func (uow *InMemoryUOW) Commit() error {
	uow.commited = true
	return nil
}

func (uow *InMemoryUOW) Rollback() {
	uow.rollbacked = true
}

type GormSqlRepository struct {
	DB *gorm.DB
}

func (r *GormSqlRepository) Add(team *domain.Team) error {
	result := r.DB.Create(team)

	return result.Error // TODO: implement own errors
}

func (r *GormSqlRepository) Get(uuid string) (*domain.Team, error) {
	team := &domain.Team{}
	result := r.DB.Model(team).Preload("Members").First(team, "uuid = ?", uuid)

	return team, result.Error
}

type GormSqlUnitOfWork struct {
	db    *gorm.DB
	tx    *gorm.DB
	teams *GormSqlRepository
}

func (uow *GormSqlUnitOfWork) Begin() error {
	uow.tx = uow.db.Begin()
	uow.teams = &GormSqlRepository{uow.db}

	return uow.tx.Error
}

func (uow *GormSqlUnitOfWork) Teams() service_layer.Repository {
	return uow.teams
}

func (uow *GormSqlUnitOfWork) Commit() error {
	uow.tx.Commit()

	return uow.tx.Error
}

func (uow *GormSqlUnitOfWork) Rollback() {
	uow.tx.Rollback()
}

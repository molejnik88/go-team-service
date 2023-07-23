package adapters

import (
	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
	"gorm.io/gorm"
)

type GormSqlRepository struct {
	DB *gorm.DB
}

func (r *GormSqlRepository) Add(team *domain.Team) error {
	return r.DB.Create(team).Error
}

func (r *GormSqlRepository) Get(uuid string) (*domain.Team, error) {
	team := &domain.Team{}
	result := r.DB.Model(team).Preload("Members").First(team, "uuid = ?", uuid)

	return team, result.Error
}

func (r *GormSqlRepository) Update(team *domain.Team) error {
	return r.DB.Save(team).Error
}

type GormSqlUnitOfWork struct {
	db    *gorm.DB
	tx    *gorm.DB
	teams *GormSqlRepository
}

func (uow *GormSqlUnitOfWork) Begin() error {
	uow.tx = uow.db.Begin()
	uow.teams = &GormSqlRepository{uow.tx}

	return uow.tx.Error
}

func (uow *GormSqlUnitOfWork) Teams() service_layer.Repository {
	return uow.teams
}

func (uow *GormSqlUnitOfWork) Commit() error {
	return uow.tx.Commit().Error
}

func (uow *GormSqlUnitOfWork) Rollback() {
	if r := recover(); r != nil {
		uow.tx.Rollback()
	}
}

func NewSqlUnitOfWork(db *gorm.DB) *GormSqlUnitOfWork {
	return &GormSqlUnitOfWork{db: db}
}

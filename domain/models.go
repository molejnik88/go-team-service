package domain

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	UUID        string `gorm:"primaryKey"`
	Name        string
	Description string
	Members     []TeamMember `gorm:"foreignKey:TeamUUID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (team *Team) AddMember(m TeamMember) error {
	// TODO: Limit number of members
	team.Members = append(team.Members, m)

	return nil
}

type TeamMember struct {
	UUID      string `gorm:"primaryKey"`
	TeamUUID  string
	Email     string
	IsAdmin   bool
	IsOwner   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

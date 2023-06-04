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

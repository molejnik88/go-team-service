package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     int
	DbSSLMode  string
}

func (c *Config) DbDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.DbHost, c.DbUser, c.DbPassword, c.DbName, c.DbPort, c.DbSSLMode,
	)
}

type Service struct {
	Router *gin.Engine
	config *Config
}

func (srv *Service) setupRouter() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/teams", createTeam)
	r.GET("/teams/:uuid", fetchTeam)

	srv.Router = r
}

func (srv *Service) connectDatabase() {
	dsn := srv.config.DbDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = db
}

func (srv *Service) Setup() {
	srv.setupRouter()
	srv.connectDatabase()
}

func (srv *Service) Run() {
	srv.Router.Run()
}

func NewServiceWithConfig(c Config) *Service {
	return &Service{config: &c}
}

func main() {
	conf := Config{
		DbHost:     "postgres",
		DbUser:     "manager",
		DbPassword: "password",
		DbName:     "team_app",
		DbPort:     5432,
		DbSSLMode:  "disable",
	}
	srv := NewServiceWithConfig(conf)
	srv.Setup()
	srv.Run()
}

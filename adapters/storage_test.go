package adapters

import (
	"testing"

	"github.com/google/uuid"
	"github.com/molejnik88/go-team-service/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSqlRepository(t *testing.T) {
	t.Run("Dummy test", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.Nil(t, err)

		err = db.AutoMigrate(&domain.Team{}, &domain.TeamMember{})
		assert.Nil(t, err)

		team := &domain.Team{
			UUID:        uuid.NewString(),
			Name:        "Test sql",
			Description: "Test create team",
			Members: []domain.TeamMember{
				{
					UUID:    uuid.NewString(),
					Email:   "fake@example.com",
					IsAdmin: true,
					IsOwner: true,
				},
			},
		}

		result := db.Create(team)
		assert.Nil(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	})

	t.Run("Create and Get - GormSqlRepository", func(t *testing.T) {
		// TODO: move db setup to test setup; check testing.Main
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.Nil(t, err)

		err = db.AutoMigrate(&domain.Team{}, &domain.TeamMember{})
		assert.Nil(t, err)

		repo := GormSqlRepository{db.Session(&gorm.Session{FullSaveAssociations: true})}

		createTeam := &domain.Team{
			UUID:        uuid.NewString(),
			Name:        "Test repo",
			Description: "Test create team",
			Members: []domain.TeamMember{
				{
					UUID:    uuid.NewString(),
					Email:   "fake@example.com",
					IsAdmin: true,
					IsOwner: true,
				},
			},
		}
		err = repo.Add(createTeam)
		assert.Nil(t, err)

		fetchTeam, err := repo.Get(createTeam.UUID)
		assert.Nil(t, err)

		assert.Equal(t, createTeam.UUID, fetchTeam.UUID)
		assert.Equal(t, createTeam.Name, fetchTeam.Name)
		assert.Equal(t, createTeam.Description, fetchTeam.Description)
		assert.Equal(t, 1, len(fetchTeam.Members))
		assert.Equal(t, createTeam.Members[0].UUID, fetchTeam.Members[0].UUID)
		assert.Equal(t, createTeam.Members[0].Email, fetchTeam.Members[0].Email)
		assert.Equal(t, createTeam.Members[0].IsAdmin, fetchTeam.Members[0].IsAdmin)
		assert.Equal(t, createTeam.Members[0].IsOwner, fetchTeam.Members[0].IsOwner)
	})
}

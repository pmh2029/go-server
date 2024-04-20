package migrations

import (
	"go-server/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		entities.User{},
	)

	return err
}

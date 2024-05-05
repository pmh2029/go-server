package migrations

import (
	"go-server/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.SetupJoinTable(&entities.Place{}, "Categories", &entities.PlaceCategory{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		entities.User{},
		entities.Banner{},
		entities.Category{},
		entities.Place{},
	)

	return err
}

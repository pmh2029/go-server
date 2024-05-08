package interfaces

import (
	"go-server/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

type PlaceCategoryRepository interface {
	BatchCreateWithTx(
		tx *gorm.DB,
		placeCategories []entities.PlaceCategory,
	) error
	DeleteByConditionsWithTx(
		tx *gorm.DB,
		conditions map[string]interface{},
	) error
}

type PlaceCategoryUsecase interface{}

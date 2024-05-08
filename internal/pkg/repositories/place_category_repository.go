package repositories

import (
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type placeCategoryRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPlaceCategoryRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) interfaces.PlaceCategoryRepository {
	return &placeCategoryRepository{
		db,
		logger,
	}
}

func (r *placeCategoryRepository) BatchCreateWithTx(
	tx *gorm.DB,
	placeCategories []entities.PlaceCategory,
) error {
	return tx.Create(&placeCategories).Error
}

func (r *placeCategoryRepository) DeleteByConditionsWithTx(
	tx *gorm.DB,
	conditions map[string]interface{},
) error {
	return tx.Where(conditions).Delete(&entities.PlaceCategory{}).Error
}

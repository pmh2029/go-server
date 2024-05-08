package repositories

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/database"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type placeRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPlaceRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) interfaces.PlaceRepository {
	return &placeRepository{
		db,
		logger,
	}
}

func (r *placeRepository) CreateWithTx(
	tx *gorm.DB,
	place entities.Place,
) (entities.Place, error) {
	err := tx.Create(&place).Error
	return place, err
}

func (r *placeRepository) FindListPaginate(
	ctx context.Context,
	pageData map[string]int,
	conditions map[string]interface{},
) ([]entities.Place, int64, error) {
	cdb := r.db.WithContext(ctx)

	var places []entities.Place
	var count int64

	countBuilder := cdb.Model(entities.Place{})
	queryBuilder := cdb.Scopes(database.Pagination(pageData))

	if keyword, ok := conditions["keyword"].(string); ok {
		delete(conditions, "keyword")
		queryBuilder = queryBuilder.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
		countBuilder = countBuilder.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if categoryID, ok := conditions["category_id"]; ok {
		delete(conditions, "category_id")
		queryBuilder = queryBuilder.
			Joins("JOIN place_categories ON (place_categories.place_id = places.id AND place_categories.deleted_at IS NULL)").
			Where("category_id = ?", categoryID)
		countBuilder = countBuilder.
			Joins("JOIN place_categories ON (place_categories.place_id = places.id AND place_categories.deleted_at IS NULL)").
			Where("category_id = ?", categoryID)
	}

	err := countBuilder.Where(conditions).Count(&count).Error
	if err != nil {
		return []entities.Place{}, 0, err
	}

	err = queryBuilder.Preload("Categories").Where(conditions).Order("updated_at DESC").Find(&places).Error
	if err != nil {
		return []entities.Place{}, 0, err
	}

	return places, count, nil
}

func (r *placeRepository) UpdateWithTx(
	tx *gorm.DB,
	place entities.Place,
	newPlace entities.Place,
) (entities.Place, error) {
	err := tx.Model(&place).Updates(newPlace).Error
	return place, err
}

func (r *placeRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Place, error) {
	cdb := r.db.WithContext(ctx)

	var place entities.Place
	err := cdb.Where(conditions).Take(&place).Error

	return place, err
}

func (r *placeRepository) TakeByConditionsWithPreload(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Place, error) {
	cdb := r.db.WithContext(ctx)

	var place entities.Place
	err := cdb.Preload("Categories").Where(conditions).Take(&place).Error

	return place, err
}

func (r *placeRepository) DeleteByConditionsWithTx(
	tx *gorm.DB,
	conditions map[string]interface{},
) error {
	return tx.Where(conditions).Delete(&entities.Place{}).Error
}

func (r *placeRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Place, error) {
	cdb := r.db.WithContext(ctx)

	var places []entities.Place

	queryBuilder := cdb

	if keyword, ok := conditions["keyword"].(string); ok {
		delete(conditions, "keyword")
		queryBuilder = queryBuilder.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if categoryID, ok := conditions["category_id"]; ok {
		delete(conditions, "category_id")
		queryBuilder = queryBuilder.
			Joins("JOIN place_categories ON (place_categories.place_id = places.id AND place_categories.deleted_at IS NULL)").
			Where("category_id = ?", categoryID)
	}

	err := queryBuilder.Preload("Categories").Where(conditions).Order("updated_at DESC").Find(&places).Error
	if err != nil {
		return []entities.Place{}, err
	}

	return places, nil
}

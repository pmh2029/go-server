package repositories

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewCategoryRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) interfaces.CategoryRepository {
	return &categoryRepository{
		db,
		logger,
	}
}

func (r *categoryRepository) Create(
	ctx context.Context,
	category entities.Category,
) (entities.Category, error) {
	cdb := r.db.WithContext(ctx)
	err := cdb.Create(&category).Error
	return category, err
}

func (r *categoryRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Category, error) {
	cdb := r.db.WithContext(ctx)

	var categories []entities.Category
	err := cdb.Where(conditions).Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) Update(
	ctx context.Context,
	category entities.Category,
	newBanner entities.Category,
) (entities.Category, error) {
	cdb := r.db.WithContext(ctx)
	err := cdb.Model(&category).Updates(newBanner).Error
	return category, err
}

func (r *categoryRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Category, error) {
	cdb := r.db.WithContext(ctx)

	var category entities.Category
	err := cdb.Where(conditions).Take(&category).Error
	return category, err
}

func (r *categoryRepository) DeleteByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) error {
	cdb := r.db.WithContext(ctx)
	return cdb.Where(conditions).Delete(&entities.Category{}).Error
}

func (r *categoryRepository) TakeByConditionsWithPreload(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Category, error) {
	cdb := r.db.WithContext(ctx)

	var category entities.Category
	err := cdb.Preload("Places").Where(conditions).Take(&category).Error
	return category, err
}

func (r *categoryRepository) PluckIDByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]int, error) {
	cdb := r.db.WithContext(ctx)

	var ids []int
	err := cdb.Model(&entities.Category{}).Where(conditions).Pluck("id", &ids).Error

	return ids, err
}

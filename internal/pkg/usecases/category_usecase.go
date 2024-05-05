package usecases

import (
	"context"
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	UpdateCategoryIDNotFound       = errors.New("Category not found")
	DeleteCategoryIDNotFound       = errors.New("Category not found")
	DeleteCategoryCategoryHasPlace = errors.New("Cannot delete category that has places")
)

type categoryUsecase struct {
	categoryRepo interfaces.CategoryRepository
	logger       *logrus.Logger
}

func NewCategoryRepository(
	categoryRepo interfaces.CategoryRepository,
	logger *logrus.Logger,
) interfaces.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo,
		logger,
	}
}

func (u *categoryUsecase) Create(
	ctx context.Context,
	req dtos.CreateCategoryRequestDto,
) (entities.Category, error) {
	category := entities.Category{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
	}

	category, err := u.categoryRepo.Create(ctx, category)
	return category, err
}

func (u *categoryUsecase) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Category, error) {
	categories, err := u.categoryRepo.FindByConditions(ctx, map[string]interface{}{})

	return categories, err
}

func (u *categoryUsecase) Update(
	ctx context.Context,
	req dtos.UpdateCategoryRequestDto,
	conditions map[string]interface{},
) (entities.Category, error) {
	category, err := u.categoryRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Category{}, UpdateCategoryIDNotFound
		}
		return entities.Category{}, err
	}

	updatedCategory := entities.Category{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
	}

	category, err = u.categoryRepo.Update(ctx, category, updatedCategory)
	if err != nil {
		return entities.Category{}, err
	}

	return category, nil
}

func (u *categoryUsecase) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Category, error) {
	return u.categoryRepo.TakeByConditions(ctx, conditions)
}

func (u *categoryUsecase) DeleteByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) error {
	category, err := u.categoryRepo.TakeByConditionsWithPreload(ctx, conditions)
	if err != nil {
		return DeleteCategoryIDNotFound
	}

	if len(category.Places) > 0 {
		return DeleteCategoryCategoryHasPlace
	}

	return u.categoryRepo.DeleteByConditions(ctx, conditions)
}

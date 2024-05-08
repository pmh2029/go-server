package interfaces

import (
	"context"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
)

type CategoryRepository interface {
	Create(
		ctx context.Context,
		category entities.Category,
	) (entities.Category, error)
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Category, error)
	Update(
		ctx context.Context,
		category entities.Category,
		updatedCategory entities.Category,
	) (entities.Category, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Category, error)
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) error
	TakeByConditionsWithPreload(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Category, error)
	PluckIDByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]int, error)
}

type CategoryUsecase interface {
	Create(
		ctx context.Context,
		req dtos.CreateCategoryRequestDto,
	) (entities.Category, error)
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Category, error)
	Update(
		ctx context.Context,
		req dtos.UpdateCategoryRequestDto,
		conditions map[string]interface{},
	) (entities.Category, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Category, error)
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) error
}

package interfaces

import (
	"context"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

type PlaceRepository interface {
	CreateWithTx(
		tx *gorm.DB,
		place entities.Place,
	) (entities.Place, error)
	FindListPaginate(
		ctx context.Context,
		pageData map[string]int,
		conditions map[string]interface{},
	) ([]entities.Place, int64, error)
	UpdateWithTx(
		tx *gorm.DB,
		place entities.Place,
		newPlace entities.Place,
	) (entities.Place, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Place, error)
	TakeByConditionsWithPreload(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Place, error)
	DeleteByConditionsWithTx(
		tx *gorm.DB,
		conditions map[string]interface{},
	) error
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Place, error)
}

type PlaceUsecase interface {
	Create(
		ctx context.Context,
		db *gorm.DB,
		req dtos.CreatePlaceRequestDto,
	) (entities.Place, error, map[string]interface{})
	FindListPaginate(
		ctx context.Context,
		pageData map[string]int,
		conditions map[string]interface{},
	) ([]entities.Place, int64, error)
	Update(
		ctx context.Context,
		db *gorm.DB,
		placeID int,
		req dtos.UpdatePlaceRequestDto,
	) (entities.Place, error, map[string]interface{})
	TakeByConditionsWithPreload(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Place, error)
	Delete(
		ctx context.Context,
		db *gorm.DB,
		conditions map[string]interface{},
	) error
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Place, error)
}

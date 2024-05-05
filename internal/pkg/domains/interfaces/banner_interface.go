package interfaces

import (
	"context"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
)

type BannerRepository interface {
	Create(
		ctx context.Context,
		banner entities.Banner,
	) (entities.Banner, error)
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Banner, error)
	Update(
		ctx context.Context,
		banner entities.Banner,
		newBanner entities.Banner,
	) (entities.Banner, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Banner, error)
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) error
}

type BannerUsecase interface {
	Create(
		ctx context.Context,
		req dtos.CreateBannerRequestDto,
	) (entities.Banner, error)
	FindByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) ([]entities.Banner, error)
	Update(
		ctx context.Context,
		req dtos.UpdateBannerRequestDto,
		conditions map[string]interface{},
	) (entities.Banner, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) (entities.Banner, error)
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]interface{},
	) error
}

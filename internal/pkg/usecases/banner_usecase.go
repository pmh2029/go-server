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
	UpdateBannerIDNotFound = errors.New("Banner not found")
	DeleteBannerIDNotFound = errors.New("Banner not found")
)

type bannerUsecase struct {
	bannerRepo interfaces.BannerRepository
	logger     *logrus.Logger
}

func NewBannerUsecase(
	bannerRepo interfaces.BannerRepository,
	logger *logrus.Logger,
) interfaces.BannerUsecase {
	return &bannerUsecase{
		bannerRepo,
		logger,
	}
}

func (u *bannerUsecase) Create(
	ctx context.Context,
	req dtos.CreateBannerRequestDto,
) (entities.Banner, error) {
	banner := entities.Banner{
		Name:  req.Name,
		Image: req.Image,
	}

	banner, err := u.bannerRepo.Create(ctx, banner)
	return banner, err
}

func (u *bannerUsecase) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Banner, error) {
	banners, err := u.bannerRepo.FindByConditions(ctx, map[string]interface{}{})

	return banners, err
}

func (u *bannerUsecase) Update(
	ctx context.Context,
	req dtos.UpdateBannerRequestDto,
	conditions map[string]interface{},
) (entities.Banner, error) {
	banner, err := u.bannerRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Banner{}, UpdateBannerIDNotFound
		}
		return entities.Banner{}, err
	}

	newBanner := entities.Banner{
		Name:  req.Name,
		Image: req.Image,
	}

	banner, err = u.bannerRepo.Update(ctx, banner, newBanner)
	if err != nil {
		return entities.Banner{}, err
	}

	return banner, nil
}

func (u *bannerUsecase) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Banner, error) {
	return u.bannerRepo.TakeByConditions(ctx, conditions)
}

func (u *bannerUsecase) DeleteByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) error {
	_, err := u.bannerRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		return DeleteBannerIDNotFound
	}

	return u.bannerRepo.DeleteByConditions(ctx, conditions)
}

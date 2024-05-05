package usecases

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
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

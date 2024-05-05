package repositories

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type bannerRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewBannerRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) interfaces.BannerRepository {
	return &bannerRepository{
		db,
		logger,
	}
}

func (r *bannerRepository) Create(
	ctx context.Context,
	banner entities.Banner,
) (entities.Banner, error) {
	cdb := r.db.WithContext(ctx)
	err := cdb.Create(&banner).Error
	return banner, err
}

func (r *bannerRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Banner, error) {
	cdb := r.db.WithContext(ctx)

	var banners []entities.Banner
	err := cdb.Where(conditions).Find(&banners).Error

	return banners, err
}

package repositories

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) interfaces.UserRepository {
	return &userRepository{
		db,
		logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	cdb := r.db.WithContext(ctx)

	err := cdb.Create(&user).Error

	return user, err
}

func (r *userRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error) {
	cdb := r.db.WithContext(ctx)

	var user entities.User
	err := cdb.Where(conditions).Take(&user).Error

	return user, err
}

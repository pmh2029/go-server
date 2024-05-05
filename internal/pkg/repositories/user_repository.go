package repositories

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"time"

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

func (r *userRepository) Update(
	ctx context.Context,
	user entities.User,
	req dtos.UpdateUserRequestDto,
) (entities.User, error) {
	cdb := r.db.WithContext(ctx)
	var birthDay time.Time

	if req.BirthDay != 0 {
		birthDay = time.Unix(int64(req.BirthDay), 0)
	}
	err := cdb.Model(&user).Updates(entities.User{
		Username: req.Username,
		Contact:  req.Contact,
		BirthDay: &birthDay,
		Gender:   req.Gender,
		Avatar:   req.Avatar,
	}).Error

	return user, err
}

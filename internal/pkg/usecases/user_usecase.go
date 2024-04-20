package usecases

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/utils"

	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo interfaces.UserRepository
	logger   *logrus.Logger
}

func NewUserUsecase(
	userRepo interfaces.UserRepository,
	logger *logrus.Logger,
) interfaces.UserUsecase {
	return &userUsecase{
		userRepo,
		logger,
	}
}

func (u *userUsecase) Create(ctx context.Context, user entities.User) (entities.User, error) {
	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPass
	user, err = u.userRepo.Create(ctx, user)
	return user, err
}

func (u *userUsecase) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error) {
	user, err := u.userRepo.TakeByConditions(ctx, conditions)
	return user, err
}

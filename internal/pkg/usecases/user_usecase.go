package usecases

import (
	"context"
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/auth"
	"go-server/pkg/shared/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	EmailNotFound = errors.New("Email not found")
	WrongPassword = errors.New("Wrong password")
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

func (u *userUsecase) TakeByConditionsWithPassword(ctx context.Context, conditions map[string]interface{}, password string) (entities.User, error) {
	user, err := u.userRepo.TakeByConditions(ctx, conditions)

	checkHashPass := utils.CheckPasswordHash(password, user.Password)
	if !checkHashPass {
		return entities.User{}, errors.New("Wrong password")
	}

	return user, err
}

func (u *userUsecase) Login(
	ctx context.Context,
	req dtos.LoginRequestDto,
) (entities.User, string, error) {
	user, err := u.userRepo.TakeByConditions(ctx, map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, "", EmailNotFound
		}
		return entities.User{}, "", err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return entities.User{}, "", WrongPassword
	}

	accessToken, err := auth.GenerateHS256JWT(map[string]interface{}{
		"user_id": user.ID,
		"sub":     user.Username,
		"email":   user.Email,
	})
	if err != nil {
		return entities.User{}, "", err
	}

	return user, accessToken, nil
}

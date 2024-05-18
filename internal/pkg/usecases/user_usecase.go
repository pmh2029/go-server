package usecases

import (
	"context"
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/auth"
	"go-server/pkg/shared/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	RegisterUsernameExisted = errors.New("Username existed")
	RegisterEmailExisted    = errors.New("Email existed")

	EmailNotFound = errors.New("Email not found")
	WrongPassword = errors.New("Wrong password")

	UpdateUserIDNotFound = errors.New("User not found")
	DetailUserIDNotFound = errors.New("User not found")
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
	if err != nil {
		if strings.Contains(err.Error(), "users.uni_users_username") {
			return entities.User{}, RegisterUsernameExisted
		}
		if strings.Contains(err.Error(), "users.uni_users_email") {
			return entities.User{}, RegisterEmailExisted
		}
		return entities.User{}, err
	}
	return user, err
}

func (u *userUsecase) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error) {
	user, err := u.userRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, DetailUserIDNotFound
		}
		return entities.User{}, err
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
		"user_id":  user.ID,
		"sub":      user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
	})
	if err != nil {
		return entities.User{}, "", err
	}

	return user, accessToken, nil
}

func (u *userUsecase) Update(
	ctx context.Context,
	conditions map[string]interface{},
	req dtos.UpdateUserRequestDto,
) (entities.User, error) {
	user, err := u.userRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, UpdateUserIDNotFound
		}
		return entities.User{}, err
	}

	user, err = u.userRepo.Update(ctx, user, req)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

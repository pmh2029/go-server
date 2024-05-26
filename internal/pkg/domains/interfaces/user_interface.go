package interfaces

import (
	"context"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error)
	Update(
		ctx context.Context,
		user entities.User,
		req dtos.UpdateUserRequestDto,
	) (entities.User, error)
}

type UserUsecase interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error)
	Login(
		ctx context.Context,
		req dtos.LoginRequestDto,
	) (entities.User, string, string, error)
	Update(
		ctx context.Context,
		conditions map[string]interface{},
		req dtos.UpdateUserRequestDto,
	) (entities.User, error)
}

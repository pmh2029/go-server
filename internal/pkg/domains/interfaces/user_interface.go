package interfaces

import (
	"context"
	"go-server/internal/pkg/domains/models/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error)
}

type UserUsecase interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entities.User, error)
}

package application

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/domain/repository"
)

type userApp struct {
	user repository.UserRepository
}

// userApp implements the UserInterface
var _ UserInterface = &userApp{}

type UserInterface interface {
	ListUser(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
}

func (u *userApp) ListUser(ctx context.Context) ([]entity.User, error) {
	return u.user.ListUser(ctx)
}

func (u *userApp) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return u.user.GetUserByID(ctx, id)
}

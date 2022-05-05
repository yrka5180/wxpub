package application

import (
	"context"
	"public-platform-manager/internal/domain/repository"
)

type userApp struct {
	user repository.UserRepository
}

// userApp implements the UserInterface
var _ UserInterface = &userApp{}

type UserInterface interface {
	ListUser(ctx context.Context) ([]interface{}, error)
	GetUser(ctx context.Context) (interface{}, error)
}

func (u *userApp) ListUser(ctx context.Context) ([]interface{}, error) {
	return u.user.ListUser(ctx)
}

func (u *userApp) GetUser(ctx context.Context) (interface{}, error) {
	return u.user.GetUser(ctx)
}

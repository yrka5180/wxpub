package repository

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/infrastructure/persistence"
)

type UserRepository struct {
	user *persistence.UserRepo
}

func NewUserRepository(user *persistence.UserRepo) *UserRepository {
	return &UserRepository{
		user: user,
	}
}

func (a *UserRepository) ListUser(ctx context.Context) ([]entity.User, error) {
	return a.user.ListUser(ctx)
}

func (a *UserRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return a.user.GetUserByID(ctx, id)
}

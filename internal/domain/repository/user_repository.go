package repository

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
)

type UserRepository struct {
	user *persistence.UserRepo
}

var defaultUserRepository = &UserRepository{}

func NewUserRepository(user *persistence.UserRepo) {
	if defaultUserRepository.user == nil {
		defaultUserRepository.user = user
	}
}

func DefaultUserRepository() *UserRepository {
	return defaultUserRepository
}

func (a *UserRepository) ListUser(ctx context.Context) ([]entity.User, error) {
	return a.user.ListUser(ctx)
}

func (a *UserRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return a.user.GetUserByID(ctx, id)
}

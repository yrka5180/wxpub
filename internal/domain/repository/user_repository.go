package repository

import "context"

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (a *UserRepository) ListUser(ctx context.Context) ([]interface{}, error) {
	return nil, nil
}

func (a *UserRepository) GetUser(ctx context.Context) (interface{}, error) {
	return nil, nil
}

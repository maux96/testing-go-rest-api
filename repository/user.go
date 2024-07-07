package repository

import (
	"context"
	"my_rest_api/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id int64) (*models.User, error)
	Close() error
}

var usedRepository UserRepository

func SetRepository(repo UserRepository) {
	usedRepository = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return usedRepository.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id int64) (*models.User, error) {
	return usedRepository.GetUserById(ctx, id)
}

func Close() error {
	return usedRepository.Close()
}

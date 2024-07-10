package repository

import (
	"context"
	"my_rest_api/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

  InsertPost(ctx context.Context, post *models.Post) error

	Close() error
}

var usedRepository Repository

func SetRepository(repo Repository) {
	usedRepository = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return usedRepository.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return usedRepository.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return usedRepository.GetUserByEmail(ctx, email)
}


func InsertPost(ctx context.Context, post *models.Post) error {
	return usedRepository.InsertPost(ctx, post)
}

func Close() error {
	return usedRepository.Close()
}

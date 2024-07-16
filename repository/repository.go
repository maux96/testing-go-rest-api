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
  GetPostById(ctx context.Context, post_id string) (*models.Post, error) 
  UpdatePost(ctx context.Context, post_id string, post *models.Post) error 
  DeletePostById(ctx context.Context, post_id string, user_id string) error 

  ListPosts(ctx context.Context, page int64) ([]*models.Post, error)
  ListPostsByUser(ctx context.Context, user_id string, page int64) ([]*models.Post, error)

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

func GetPostById(ctx context.Context, post_id string) (*models.Post, error) {
  return usedRepository.GetPostById(ctx, post_id)
} 

func UpdatePost(ctx context.Context, post_id string, post *models.Post) error  {
  return usedRepository.UpdatePost(ctx, post_id, post)
}
func DeletePostById(ctx context.Context, post_id string, user_id string) error {
  return usedRepository.DeletePostById(ctx, post_id, user_id)
}

func ListPosts(ctx context.Context, page int64) ([]*models.Post, error) {
  return usedRepository.ListPosts(ctx, page)
}
func ListPostsByUser(ctx context.Context, user_id string, page int64) ([]*models.Post, error) {
  return usedRepository.ListPostsByUser(ctx, user_id, page)
}


func Close() error {
	return usedRepository.Close()
}

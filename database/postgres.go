package database

import (
	"context"
	"database/sql"
	"log"
	"my_rest_api/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (pr *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := pr.db.ExecContext(
		ctx,
		"INSERT INTO users (id, email, hashed_password) VALUES ($1, $2, $3);",
		user.Id,
		user.Email,
		user.HashedPassword,
	)

	return err
}

func (pr *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, email, hashed_password FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.HashedPassword); err == nil {
			return &user, nil
		}
	}

	return nil, rows.Err()
}

func (pr *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, email, hashed_password FROM users WHERE email = $1;", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.HashedPassword); err == nil {
			return &user, nil
		}
	}

	return nil, rows.Err()
}


func (pr *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := pr.db.ExecContext(
		ctx,
		"INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3);",
		post.Id,
		post.PostContent,
    post.UserId,
	)

	return err
}

func (pr *PostgresRepository) Close() error {
	return pr.db.Close()
}

func (pr *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	var post models.Post
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			return &post, nil
		}
	}

	return nil, rows.Err()
}

func (pr *PostgresRepository) UpdatePost(ctx context.Context, post_id string,post *models.Post) error  {
  _, err := pr.db.ExecContext(
		ctx,
		"UPDATE posts SET post_content = $1 WHERE id = $2 AND user_id = $3 ",
		post.PostContent,
    post_id,
    post.UserId,
	)

	return err


}
func (pr *PostgresRepository) DeletePostById(ctx context.Context, post_id string, user_id string) error {
  _, err := pr.db.ExecContext(
		ctx,
		"DELETE FROM posts WHERE id = $1 AND user_id = $2",
    post_id,
    user_id,
	)
  return err 
}

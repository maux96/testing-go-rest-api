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

func (pr *PostgresRepository) Close() error {
	return pr.db.Close()
}

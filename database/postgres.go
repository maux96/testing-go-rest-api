package database

import (
	"context"
	"database/sql"
	"my_rest_api/models"
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
		"INSERT INTO users (email) VALUES ($1, $2)",
		user.Email,
	)

	return err
}

func (pr *PostgresRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(user.Id, user.Email); err != nil {
			return &user, nil
		}
	}

	return nil, rows.Err()
}

func (pr *PostgresRepository) Close() error {
	return pr.db.Close()
}

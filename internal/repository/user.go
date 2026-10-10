package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/models"
)

func RegisterUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO users (email, password)
	VALUES ($1, $2)
	RETURNING id, email, created_at, updated_at
	`

	var createdUser models.User

	if err := pool.QueryRow(ctx, query, user.Email, user.HashedPassword).Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, email, password, created_at, updated_at
	FROM users 
	WHERE email = $1
	`

	var user models.User

	if err := pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserById(pool *pgxpool.Pool, id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, email, password, created_at, updated_at
	FROM users 
	WHERE id = $1
	`

	var user models.User

	if err := pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/models"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO todos (title, completed)
	VALUES ($1, $2)
	RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo

	err := pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, title, completed, created_at, updated_at 
	FROM todos
	ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo

	for rows.Next() {
		var todo models.Todo

		if err := rows.Scan(
			&todo.Id,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, title, completed, created_at, updated_at 
	FROM todos 
	WHERE id = $1
	`

	var todo models.Todo

	if err := pool.QueryRow(ctx, query, id).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoById(pool *pgxpool.Pool, id int, newTitle string, newCompleted bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE todos
	SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo

	if err := pool.QueryRow(ctx, query, newTitle, newCompleted, id).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodoById(pool *pgxpool.Pool, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM todos 
	WHERE id = $1
	`

	tag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

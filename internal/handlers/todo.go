package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/repository"
)

type createTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type updateTodoInput struct {
	Title     string `json:"title"`
	Completed *bool  `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInput createTodoInput
		err := ctx.ShouldBindJSON(&userInput)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdTodo, err := repository.CreateTodo(
			pool,
			userInput.Title,
			userInput.Completed,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, createdTodo)
	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		todos, err := repository.GetAllTodos(pool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		todo, err := repository.GetTodoById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		var userInput updateTodoInput
		err = ctx.ShouldBindJSON(&userInput)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userInput.Title == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
			return
		}

		if userInput.Completed == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "completed cannot be empty"})
			return
		}

		updatedTodo, err := repository.UpdateTodoById(
			pool,
			id,
			userInput.Title,
			*userInput.Completed,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, updatedTodo)
	}
}

func DeleteTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		if err = repository.DeleteTodoById(pool, id); err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

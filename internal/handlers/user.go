package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/models"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData registerRequest

		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if password is at least 8 characters.
		if len(requestData.Password) < 8 {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "password must be  8 to 72 characters long"},
			)
			return
		}

		// Hash password.
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(requestData.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			if err == bcrypt.ErrPasswordTooLong {
				ctx.JSON(
					http.StatusBadRequest,
					gin.H{"error": "password must be 8 to 72 characters long"},
				)
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user := &models.User{
			Email:          requestData.Email,
			HashedPassword: string(hashedPassword),
		}

		createdUser, err := repository.RegisterUser(pool, user)
		if err != nil {
			if err.Error() != "" {
				ctx.JSON(
					http.StatusBadRequest,
					gin.H{"error": "email already registered"},
				)
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, createdUser)

	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/config"
	"github.com/rgmaz/golang-gin-postgres-todo-rest-api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer pool.Close()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Todo API is running",
			"status":  "success",
		})
	})

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}

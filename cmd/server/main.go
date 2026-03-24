package main

import (
	"follow-system/internal/config"
	"follow-system/internal/db"
	"follow-system/internal/handlers"
	"follow-system/internal/repository"
	"follow-system/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	//Config
	cfg := config.Load()

	//database connection
	dbConn, err := db.NewPostgres(cfg.DBUrl)
	if err != nil {
		log.Fatal("Error connecting db", err)
	}

	repo := repository.NewFollowRepo(dbConn)
	svc := services.NewFollowService(repo)
	handler := handlers.NewFollowHandler(svc)

	// Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/follow", handler.FollowUser)

	log.Printf("Server running on %v", cfg.Port)

	r.Run(cfg.Port)
}

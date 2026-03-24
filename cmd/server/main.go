package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Println("Server running on :8080")

	r.Run(":8080")
}

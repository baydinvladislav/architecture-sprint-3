package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	if err := router.Run(":8082"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

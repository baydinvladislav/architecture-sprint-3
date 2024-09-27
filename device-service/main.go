package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to DeviceService!"})
	})

	err := r.Run(":8081")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

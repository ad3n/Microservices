package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ad3n/microservices/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {

	r := gin.Default()

	r.GET("/hello", middlewares.ValidateToken(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Service 2",
		})
	})

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}

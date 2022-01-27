package main

import (
	"encoding/json"
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
		var client = &http.Client{}

		request, err := http.NewRequest(http.MethodGet, "http://service2:5555/hello", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		request.Header.Set("X-User-ID", c.Request.Header.Get("X-User-ID"))
		request.Header.Set("X-User-Email", c.Request.Header.Get("X-User-Email"))
		request.Header.Set("X-Request-ID", c.Request.Header.Get("X-Request-ID"))

		response, err := client.Do(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}
		defer response.Body.Close()

		var data struct {
			Message string `json:"message"`
		}

		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Service 3 called service 2 with response: %s", data.Message)})
	})

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}

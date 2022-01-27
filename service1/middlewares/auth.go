package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int64
	Email string
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User{}
		user.ID, _ = strconv.ParseInt(c.Request.Header.Get("X-User-ID"), 10, 64)
		user.Email = c.Request.Header.Get("X-User-Email")

		if user.ID <= 0 || user.Email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
			})

			return
		}

		c.Set("user-id", user.ID)
		c.Set("user-email", user.Email)

		c.Next()
	}
}

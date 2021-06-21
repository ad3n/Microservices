package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ad3n/microservices/configs"
	"github.com/ad3n/microservices/models"
	"github.com/ad3n/microservices/repositories"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Repository *repositories.UserRepository
}

func (l Login) Auth(c *gin.Context) {
	model := models.User{}
	request := models.User{}

	err := c.BindJSON(&request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request not valid",
		})

		return
	}

	model.Email = request.Email
	l.Repository.FindByEmail(&model)

	err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(request.Password))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "username and Password not match",
		})

		return
	}

	var claims = jwt.MapClaims{}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims["email"] = model.Email
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(configs.Env.AppSignKey))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "username and Password not match",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func (l Login) Validate(c *gin.Context) {
	authHeader := strings.TrimSpace(c.Request.Header.Get("Authorization"))
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "token is missing",
		})

		return
	}

	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method != jwt.GetSigningMethod("HS256") {
			return nil, fmt.Errorf("token is not valid")
		}

		return []byte(configs.Env.AppSignKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})

		return
	}

	model := models.User{}
	value, _ := token.Claims.(jwt.MapClaims)
	model.Email, _ = value["email"].(string)

	l.Repository.FindByEmail(&model)
	if model.ID < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not found",
		})

		return
	}

	model.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"data": map[string]string{
			"id":    strconv.Itoa(model.ID),
			"email": model.Email,
		},
	})
}

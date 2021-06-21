package controllers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/ad3n/microservices/models"
	"github.com/ad3n/microservices/repositories"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repository *repositories.UserRepository
}

func (u User) Create(c *gin.Context) {
	var model = models.User{}

	c.BindJSON(&model)

	var regex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !regex.MatchString(model.Email) {
		c.JSON(http.StatusBadRequest, "Email not valid")

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.MinCost)

	model.Password = string(hash)

	u.Repository.Save(&model)

	model.Password = ""

	c.JSON(http.StatusOK, model)
}

func (u User) Update(c *gin.Context) {
	var model = models.User{}
	var request = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})

		return
	}

	c.BindJSON(&request)
	model.Email = request.Email

	u.Repository.Save(&model)

	c.JSON(http.StatusOK, model)
}

func (u User) Delete(c *gin.Context) {
	var model = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})

		return
	}

	u.Repository.Storage.Delete(&model)

	c.JSON(http.StatusNoContent, "")
}

func (u User) GetAll(c *gin.Context) {
	var models = []models.User{}

	u.Repository.All(&models)

	for i, m := range models {
		m.Password = ""
		models[i] = m
	}

	c.JSON(http.StatusOK, models)
}

func (u User) Get(c *gin.Context) {
	var model = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})

		return
	}

	model.Password = ""

	c.JSON(http.StatusOK, model)
}

package repositories

import (
	"github.com/ad3n/microservices/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	Storage *gorm.DB
}

func (r *UserRepository) Save(user *models.User) {
	r.Storage.Save(user)
}

func (r *UserRepository) Find(user *models.User) {
	r.Storage.First(user)
}

func (r *UserRepository) FindByEmail(user *models.User) {
	r.Storage.Where("email = ?", user.Email).First(user)
}

func (r *UserRepository) All(user *[]models.User) {
	r.Storage.Find(user)
}

func (r *UserRepository) Remove(user *models.User) {
	r.Storage.Save(&user)
}

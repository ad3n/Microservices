package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dB *gorm.DB

func Connect() *gorm.DB {
	if dB != nil {
		return dB
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Env.DbUser, Env.DbPassword, Env.DbHost, Env.DbPort, Env.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	dB = db

	return dB
}

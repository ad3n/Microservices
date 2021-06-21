package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Env env

type env struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	DbPort     int
	AppPort    int
	AppSignKey string
}

func LoadEnv() {
	godotenv.Load()

	Env.DbUser = os.Getenv("DB_USER")
	Env.DbPassword = os.Getenv("DB_PASSWORD")
	Env.DbName = os.Getenv("DB_NAME")
	Env.DbHost = os.Getenv("DB_HOST")
	Env.DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	Env.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	Env.AppSignKey = os.Getenv("APP_SIGN_KEY")
}

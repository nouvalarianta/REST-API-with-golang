package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error when loading file configuration" + err.Error())
	}

	expInt, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
			Asset: os.Getenv("SERVER_ASSET_URL"),
		},
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Jwt: Jwt{
			Key: os.Getenv("JWT_KEY"),
			Exp: expInt,
		},
		Storage: Storage{
			BasePath: os.Getenv("STORAGE_PATH "),
		},
	}
}

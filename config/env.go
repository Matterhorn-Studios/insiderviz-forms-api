package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return os.Getenv(key)
}

package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, using system environment variables")
	}
	log.Println("environment variables loaded successfully")
	return nil
}

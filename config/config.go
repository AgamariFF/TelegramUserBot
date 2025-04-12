package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Gmail    string
	Password string
}

func LoadConfig() Config {

	if err := godotenv.Load("./config/.env"); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	config := Config{
		Gmail:    getEnv("gmail", "="),
		Password: getEnv("password", "="),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

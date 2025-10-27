package config

// its only work is to create a function that would be used to load env variables
// for this we use the godotenv package and we then can call that config file into the main.go
// in main.go then we will use the env variables as needed.

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

// Load .env file on package initialization
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname"),
	}
}

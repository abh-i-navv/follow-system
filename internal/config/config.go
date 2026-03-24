package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string
}

func Load() Config {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	return Config{
		Port:  ":" + getEnv("PORT", "8080"),
		DBUrl: getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

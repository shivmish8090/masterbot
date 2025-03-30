package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Token         string
	StartImageUrl string
)

func init() {
	godotenv.Load()

	Token = Getenv("TOKEN", "")
	if Token == "" {
		panic("TOKEN environment variable is empty")
	}
	StartImageUrl = Getenv("START_IMG_URL", "https://graph.org/file/f3c8291963a053ac18536-3558d077ad80845bd7.jpg")
}

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

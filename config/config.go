package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Token      string
	StartImage string
)

func init() {
	godotenv.Load()

	Token = Getenv("TOKEN", "8050656956:AAGtdazbq1CfUg1Ok1h4QR9rU023ZJf7cso")
	if Token == "" {
		panic("TOKEN environment variable is empty")
	}
	StartImage = Getenv("START_IMG_URL", "https://telegra.ph/file/fef61e95cc35da109f900-e3978f1e3eaad29dee.jpg")
}

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

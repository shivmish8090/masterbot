package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Token        string
	StartImageID string
)

func init() {
	godotenv.Load()

	Token = Getenv("TOKEN", "8050656956:AAGtdazbq1CfUg1Ok1h4QR9rU023ZJf7cso")
	if Token == "" {
		panic("TOKEN environment variable is empty")
	}
	StartImageID = Getenv("START_IMG_URL", "AgACAgUAAx0Cf-yYdwABASe0Z-jiV0XQqdx2z9baWpZt-9a_r0IAAq3OMRu1GhFXGqclfejWhCUACAEAAwIAA3kABx4E")
}

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

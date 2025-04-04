package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	LoggerId   int64
	MongoUri   string
	OwnerId    int64
	StartImage string
	Token      string
)

func init() {
	godotenv.Load()
	Token = Getenv("TOKEN", "8050656956:AAGtdazbq1CfUg1Ok1h4QR9rU023ZJf7cso")
	StartImage = Getenv(
		"START_IMG_URL",
		"https://telegra.ph/file/ba238ec5e542d8754cea7-dc1786aa23ae1224f2.jpg",
	)
	LoggerId = GetenvInt64("LOGGER_ID", "-1002647107199")
	MongoUri = Getenv("MONGO_DB_URI", "")
	OwnerId = GetenvInt64("OWNER_ID", "7706682472")

	if Token == "" {
		log.panic("TOKEN environment variable is empty")
	}
	if MongoUri == "" {
		log.panic("MONGO_DB_URI environment variable is empty")
	}
}

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetenvInt64(key, defaultValue string) int64 {
	value := Getenv(key, defaultValue)
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error converting %s: %v", key, err))
	}
	return intValue
}

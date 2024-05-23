package helper

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	LoadEnv()
}

func LoadEnv() {
	currentWorkDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = godotenv.Load(currentWorkDirectory + "/config/.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func DecodeString(encodedStr string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedStr)
}

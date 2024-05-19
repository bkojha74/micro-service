package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client         *mongo.Client
	UserCollection *mongo.Collection
)

func Init() {
	LoadEnv()
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))
	return hashedPassword
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

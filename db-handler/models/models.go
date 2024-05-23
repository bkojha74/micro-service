package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bkojha74/micro-service/db-handler/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	SecretKey string `json:"secret" bson:"secret"`
	Role      string `json:"role" bson:"role"`
}

type UserModel interface {
	CreateUser(user User) error
	ReadUser(username string) (User, error)
	UpdateUser(user User) error
	DeleteUser(username string) error
}

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opt := options.Client().ApplyURI(helper.GetEnv("MONGO_URI"))

	var err error
	helper.Client, err = mongo.Connect(ctx, opt)
	if err != nil {
		fmt.Println("Client connection failed")
	}

	helper.UserCollection = helper.Client.Database("filehandler").Collection("users")

	fmt.Println("Database connected.\nReady to server APIs")
}

type MongoUserModel struct{}

func (m *MongoUserModel) CreateUser(user User) error {
	user.Password = helper.HashString(user.Password)
	user.SecretKey = helper.EncodeString(user.SecretKey)

	_, err := helper.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (m *MongoUserModel) ReadUser(username string) (User, error) {
	var user User
	filter := bson.M{"username": username}
	err := helper.UserCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println("Error reading user:", err)
		return user, err
	}
	return user, nil
}

func (m *MongoUserModel) UpdateUser(user User) error {
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": bson.M{"password": helper.HashString(user.Password)}}
	_, err := helper.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

func (m *MongoUserModel) DeleteUser(username string) error {
	filter := bson.M{"username": username}
	_, err := helper.UserCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

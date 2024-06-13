package controller

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bkojha74/micro-service/auth-handler/models"
	proto "github.com/bkojha74/micro-service/auth-handler/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ExampleClient

type Server struct {
	proto.UnimplementedExampleServer
}

func (s *Server) GetUserInfo(ctx context.Context, in *proto.DBRequest) (*proto.DBResponse, error) {
	return &proto.DBResponse{}, errors.New("")
}

func getUserInfo(user string) models.UserwithError {
	conn, err := grpc.Dial("db-handler:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	client = proto.NewExampleClient(conn)

	req := proto.DBRequest{
		Username: user,
	}

	fmt.Println("Requesting to get User Info from db-handler")

	res, err := client.GetUserInfo(context.TODO(), &req)
	if err != nil {
		fmt.Println("Received error from db-handler:", err)
		return models.UserwithError{
			Err: err,
		}
	}

	fmt.Printf("Received response from db-handler:%+v\n", res)

	ret := models.UserwithError{
		User: models.User{
			Username:  res.Username,
			Password:  res.UserInfo.Password,
			SecretKey: res.UserInfo.Secret,
			Role:      res.UserInfo.Role,
		},
	}

	return ret
}

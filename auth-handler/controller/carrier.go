package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bkojha74/micro-service/auth-handler/models"
	proto "github.com/bkojha74/micro-service/auth-handler/protoc"
	"google.golang.org/grpc"
)

var client proto.ExampleClient

type Server struct {
	proto.UnimplementedExampleServer
}

func (s *Server) GetUserInfo(ctx context.Context, in *proto.DBRequest) (*proto.DBResponse, error) {
	return &proto.DBResponse{}, nil
}

func getUserInfo(user string) models.UserwithError {
	conn, err := grpc.Dial("db-handler:8084", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to db-handler: %v", err)
	}
	defer conn.Close()

	client = proto.NewExampleClient(conn)

	req := proto.DBRequest{
		Username: user,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("Requesting to get User Info from db-handler")

	res, err := client.GetUserInfo(ctx, &req)
	if err != nil {
		fmt.Println("Received error from db-handler:", err)
		return models.UserwithError{
			Err: err,
		}
	}

	fmt.Println("Received response from db-handler:", res)

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

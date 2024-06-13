package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/bkojha74/micro-service/db-handler/models"

	proto "github.com/bkojha74/micro-service/db-handler/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedExampleServer
}

func (s *Server) GetUserInfo(ctx context.Context, in *proto.DBRequest) (*proto.DBResponse, error) {
	fmt.Println("Request received to get User Info from auth-handler")
	usr := models.MongoUserModel{}
	usrInfo, err := usr.ReadUser(in.Username)
	if err != nil {
		return &proto.DBResponse{}, err
	}

	fmt.Println("2. Got UserInfo", usrInfo)

	return &proto.DBResponse{
		Username: usrInfo.Username,
		UserInfo: &proto.UserInfo{
			Password: usrInfo.Password,
			Secret:   usrInfo.SecretKey,
			Role:     usrInfo.Role,
		},
	}, errors.New("")
}

func StartGRPCServer() {
	listener, tcpErr := net.Listen("tcp", ":8084")
	if tcpErr != nil {
		log.Panic(tcpErr)
	}
	srv := grpc.NewServer()
	proto.RegisterExampleServer(srv, &Server{})
	reflection.Register(srv)

	log.Println("gRPC server is running on port 8084")

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

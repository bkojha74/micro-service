package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/bkojha74/micro-service/auth-handler/docs"
	"github.com/bkojha74/micro-service/auth-handler/helper"
	proto "github.com/bkojha74/micro-service/auth-handler/protoc"
	route "github.com/bkojha74/micro-service/auth-handler/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title Authentication Management API
// @version 1.0
// @description This is a sample server for Authentication management.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /
func main() {
	helper.Init()

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(err)
	}
	client := proto.NewExampleClient(conn)
	client.GetUserInfo(context.TODO(), &proto.DBRequest{})

	router := route.GetRouter()

	// Start the server
	log.Fatal(http.ListenAndServe(":8081", router))
}

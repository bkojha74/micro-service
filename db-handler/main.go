package main

import (
	"log"
	"net/http"

	"github.com/bkojha74/micro-service/db-handler/controller"
	_ "github.com/bkojha74/micro-service/db-handler/docs"
	"github.com/bkojha74/micro-service/db-handler/helper"
	"github.com/bkojha74/micro-service/db-handler/models"
	route "github.com/bkojha74/micro-service/db-handler/router"
)

// @title Database Management API
// @version 1.0
// @description This is a sample server for database management.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8082
// @BasePath /
func main() {
	helper.Init()

	models.Init()

	go controller.StartGRPCServer()

	// Initialize the UserModel
	userModel := &models.MongoUserModel{}

	router := route.GetRouter(userModel)

	// Start the server
	log.Fatal(http.ListenAndServe(":8082", router))
}

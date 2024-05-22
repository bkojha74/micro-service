package main

import (
	"log"
	"net/http"

	_ "github.com/bkojha74/micro-service/file-handler/docs"
	route "github.com/bkojha74/micro-service/file-handler/router"
)

// @title File Management API
// @version 1.0
// @description This is a sample server for file management.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8083
// @BasePath /
func main() {
	router := route.GetRouter()

	// Start the server
	log.Fatal(http.ListenAndServe(":8083", router))
}

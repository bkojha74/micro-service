package router

import (
	"net/http"

	"github.com/bkojha74/micro-service/file-handler/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter() *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// Register handlers
	route.Handle("/searchdir", controller.VerifyToken(http.HandlerFunc(controller.SearchDirHandler))).Methods("GET")
	route.Handle("/file", controller.VerifyToken(http.HandlerFunc(controller.FileHandler))).Methods("POST")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return route
}

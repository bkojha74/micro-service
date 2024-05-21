package router

import (
	"github.com/bkojha74/micro-service/file-handler/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter() *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// Register handlers
	route.HandleFunc("/searchdir", controller.SearchDirHandler).Methods("GET")
	route.HandleFunc("/file", controller.FileHandler).Methods("POST")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return route
}

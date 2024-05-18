package router

import (
	"github.com/bkojha74/mocro-service/auth-handler/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter() *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// token handlers
	route.HandleFunc("/generate-token", controller.GenerateToken).Methods("POST")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return route
}

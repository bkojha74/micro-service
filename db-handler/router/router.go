package router

import (
	"github.com/bkojha74/micro-service/db-handler/controller"
	"github.com/bkojha74/micro-service/db-handler/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter(userModel models.UserModel) *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// Create a new handler instance
	handler := &controller.Handler{UserModel: userModel}

	// Crud handler
	route.HandleFunc("/users", handler.CreateUserHandler).Methods("POST")
	route.HandleFunc("/users/{username}", handler.ReadUserHandler).Methods("GET")
	route.HandleFunc("/users", handler.UpdateUserHandler).Methods("PUT")
	route.HandleFunc("/users/{username}", handler.DeleteUserHandler).Methods("DELETE")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return route
}

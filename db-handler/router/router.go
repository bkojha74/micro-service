package router

import (
	"net/http"

	"github.com/bkojha74/micro-service/db-handler/controller"
	"github.com/bkojha74/micro-service/db-handler/metrics"
	"github.com/bkojha74/micro-service/db-handler/models"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter(userModel models.UserModel) *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// Create a new handler instance
	handler := &controller.Handler{UserModel: userModel}

	// Crud handler
	route.Handle("/users", metrics.RecordMetrics("DB-Handler-Microservice", "CreateUserHandler", http.HandlerFunc(handler.CreateUserHandler))).Methods("POST")
	route.Handle("/users/{username}", metrics.RecordMetrics("DB-Handler-Microservice", "ReadUserHandler", http.HandlerFunc(handler.ReadUserHandler))).Methods("GET")
	route.Handle("/users", metrics.RecordMetrics("DB-Handler-Microservice", "UpdateUserHandler", http.HandlerFunc(handler.UpdateUserHandler))).Methods("PUT")
	route.Handle("/users/{username}", metrics.RecordMetrics("DB-Handler-Microservice", "DeleteUserHandler", http.HandlerFunc(handler.DeleteUserHandler))).Methods("DELETE")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	route.Handle("/metrics", promhttp.Handler())

	return route
}

package router

import (
	"net/http"

	"github.com/bkojha74/micro-service/auth-handler/controller"
	"github.com/bkojha74/micro-service/auth-handler/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter() *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// token handlers
	//route.HandleFunc("/generate-token", controller.GenerateToken).Methods("POST")
	//route.HandleFunc("/verify-token", controller.VerifyToken).Methods("GET")
	route.Handle("/generate-token", metrics.RecordMetrics("Auth-Handler-Microservice", "GenerateToken", http.HandlerFunc(controller.GenerateToken))).Methods("POST")
	route.Handle("/verify-token", metrics.RecordMetrics("Auth-Handler-Microservice", "VerifyToken", http.HandlerFunc(controller.VerifyToken))).Methods("GET")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	route.Handle("/metrics", promhttp.Handler())

	return route
}

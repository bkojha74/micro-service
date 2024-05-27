package router

import (
	"net/http"

	"github.com/bkojha74/micro-service/file-handler/controller"
	"github.com/bkojha74/micro-service/file-handler/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRouter() *mux.Router {
	// Create a new router
	route := mux.NewRouter()

	// Middleware chain combining VerifyToken and PrometheusMiddleware
	middlewareChain := controller.ChainMiddleware(controller.VerifyToken, metrics.PrometheusMiddleware)

	// Register handlers with middleware chain
	route.Handle("/searchdir", middlewareChain(http.HandlerFunc(controller.SearchDirHandler))).Methods("GET")
	route.Handle("/file", middlewareChain(http.HandlerFunc(controller.FileHandler))).Methods("POST")

	// Serve Swagger UI
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Add Prometheus metrics endpoint
	route.Handle("/metrics", promhttp.Handler()).Methods("GET")

	return route
}

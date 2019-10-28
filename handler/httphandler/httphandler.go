package httphandler

import (
	"net/http"
	_ "net/http/pprof" // pprof enable

	"github.com/Toshik1978/go-rest-api/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewHTTPHandler creates new http handler
func NewHTTPHandler(globals service.Globals) http.Handler {
	// Create main router and attach common middlewares
	r := mux.NewRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.Use(
		handlers.RecoveryHandler(handlers.RecoveryLogger(newRecoveryLogger(globals))),
		handlers.ProxyHeaders,
	)

	// API
	route := r.PathPrefix("/v1").Subrouter()
	route.Use(
		func(next http.Handler) http.Handler {
			return handlers.CustomLoggingHandler(nil, next, newLogFormatter(globals))
		},
	)

	apiHandler := newAPIHandler(globals)
	route.Handle("/server/status", apiHandler.ServerStatusHandler()).Methods("GET")

	return r
}

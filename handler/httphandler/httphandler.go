package httphandler

import (
	"net/http"
	_ "net/http/pprof" // pprof enable

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewHTTPHandler creates new http handler
func NewHTTPHandler(globals server.Globals, accountManager handler.AccountManager) http.Handler {
	// Create main router and attach common middlewares
	r := mux.NewRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.Use(
		handlers.RecoveryHandler(handlers.RecoveryLogger(newRecoveryLogger(globals))),
		handlers.ProxyHeaders,
	)

	// API
	route := r.PathPrefix("/api/v1").Subrouter()
	route.Use(
		func(next http.Handler) http.Handler {
			return handlers.CustomLoggingHandler(nil, next, newLogFormatter(globals))
		},
	)

	apiHandler := newAPIHandler(globals, accountManager)
	route.Handle("/server/status", apiHandler.ServerStatusHandler()).Methods("GET")

	route.Handle("/accounts", apiHandler.CreateAccountHandler()).Methods("POST")
	route.Handle("/accounts", apiHandler.GetAllAccountsHandler()).Methods("GET")
	route.Handle("/accounts/payments", apiHandler.GetAllPaymentsHandler()).Methods("GET")
	route.Handle("/accounts/{uid:[a-zA-Z0-9]+}/payments", apiHandler.CreatePaymentHandler()).Methods("POST")

	return r
}

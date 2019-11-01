package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/service/errutil"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	uidKey = "uid"
)

// apiHandler declare handler API requests
type apiHandler struct {
	logger    *zap.Logger
	buildTime string
	version   string

	accountManager handler.AccountManager
}

// newAPIHandler creates new API handler
func newAPIHandler(globals server.Globals, accountManager handler.AccountManager) *apiHandler {
	return &apiHandler{
		logger:         globals.Logger,
		buildTime:      globals.BuildTime,
		version:        globals.Version,
		accountManager: accountManager,
	}
}

// ServerStatusHandler response with status of the service and it's version
func (h *apiHandler) ServerStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.writeResponse(w, handler.ServerStatusResponse{
			IsAlive:   true,
			BuildTime: h.buildTime,
			Version:   h.version,
		})
	})
}

// CreateAccountHandler creates new account
func (h *apiHandler) CreateAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			h.fail(w, errors.New("no body detected"), http.StatusBadRequest, "CreateAccountHandler")
			return
		}

		var accountRequest handler.AccountRequest
		decoder := json.NewDecoder(r.Body)
		if h.fail(w,
			errutil.Wrap(decoder.Decode(&accountRequest), "decode failed"),
			http.StatusBadRequest, "CreateAccountHandler") {

			return
		}

		account, err := h.accountManager.AccountBuilder().
			SetUID(accountRequest.UID).
			SetCurrency(accountRequest.Currency).
			SetBalance(accountRequest.Balance).
			Build(r.Context())
		if h.fail(w,
			errutil.Wrap(err, "failed to create account"),
			http.StatusInternalServerError, "CreateAccountHandler") {

			return
		}

		w.WriteHeader(http.StatusCreated)
		h.writeResponse(w, account)
	})
}

// GetAllAccountsHandler response with all accounts
func (h *apiHandler) GetAllAccountsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accounts, err := h.accountManager.AllAccounts(r.Context())
		if h.fail(w,
			errutil.Wrap(err, "failed to get all accounts"),
			http.StatusInternalServerError, "GetAllAccountsHandler") {

			return
		}
		h.writeResponse(w, accounts)
	})
}

// GetAllPaymentsHandler response with all payments for the given account
func (h *apiHandler) GetAllPaymentsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payments, err := h.accountManager.AllPayments(r.Context())
		if h.fail(w,
			errutil.Wrap(err, "failed to get all payments"),
			http.StatusInternalServerError, "GetAllPaymentsHandler") {

			return
		}
		h.writeResponse(w, payments)
	})
}

// CreatePaymentHandler creates payment from one account to another
func (h *apiHandler) CreatePaymentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			h.fail(w, errors.New("no body detected"), http.StatusBadRequest, "CreatePaymentHandler")
			return
		}

		vars := mux.Vars(r)
		if _, ok := vars[uidKey]; !ok {
			// Theoretically it's impossible situation due to mux routing
			// But just in case...
			h.fail(w, errors.New("no payer detected"), http.StatusBadRequest, "CreatePaymentHandler")
			return
		}

		var paymentRequest handler.PaymentRequest
		decoder := json.NewDecoder(r.Body)
		if h.fail(w,
			errutil.Wrap(decoder.Decode(&paymentRequest), "decode failed"),
			http.StatusBadRequest, "CreatePaymentHandler") {

			return
		}

		payment, err := h.accountManager.PaymentBuilder().
			SetPayer(vars[uidKey]).
			SetRecipient(paymentRequest.RecipientUID).
			SetAmount(paymentRequest.Amount).
			Build(r.Context())
		if h.fail(w,
			errutil.Wrap(err, "failed to create payment"),
			http.StatusInternalServerError, "CreatePaymentHandler") {

			return
		}

		w.WriteHeader(http.StatusCreated)
		h.writeResponse(w, payment)
	})
}

// httpCode looks at error and tries to provide correct http status code or defaultCode otherwise
func (h *apiHandler) httpCode(err error, defaultCode int) int {
	var handlerError *handler.Error
	if err != nil && errors.As(err, &handlerError) {
		switch handlerError.Kind {
		case handler.ServerError:
			return http.StatusInternalServerError
		case handler.ClientError:
			return http.StatusBadRequest
		}
	}
	return defaultCode
}

// fail fails request
func (h *apiHandler) fail(w http.ResponseWriter, err error, code int, method string) bool {
	if err != nil {
		h.logger.Error("Failed to handle "+method, zap.Error(err))
		http.Error(w, err.Error(), h.httpCode(err, code))
		return true
	}
	return false
}

// writeResponse write response
func (h *apiHandler) writeResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal HTTP response", zap.Error(err), zap.Any("response", response))
		return
	}
	_, err = w.Write(payload)
	if err != nil {
		h.logger.Error("Failed to write HTTP response", zap.Error(err), zap.Any("response", response))
		return
	}
}

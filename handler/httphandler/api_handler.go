package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/service"
	"go.uber.org/zap"
)

// apiHandler declare handler API requests
type apiHandler struct {
	logger    *zap.Logger
	buildTime string
	version   string
}

// newAPIHandler creates new API handler
func newAPIHandler(globals service.Globals) *apiHandler {
	return &apiHandler{
		logger:    globals.Logger,
		buildTime: globals.BuildTime,
		version:   globals.Version,
	}
}

// ServerStatusHandler sends status of the service and it's version
func (h *apiHandler) ServerStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.writeResponse(w, handler.ServerStatusResponse{
			IsAlive:   true,
			BuildTime: h.buildTime,
			Version:   h.version,
		})
	})
}

// fail fails request
//func (h *apiHandler) fail(w http.ResponseWriter, err error, method string) bool {
//	if err != nil {
//		h.logger.Error("Failed to handle "+method, zap.Error(err))
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return true
//	}
//	return false
//}

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

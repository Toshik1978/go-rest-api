package handler

// ServerStatusResponse declare response for status request
type ServerStatusResponse struct {
	IsAlive   bool   `json:"is_alive"`
	BuildTime string `json:"build_time"`
	Version   string `json:"version"`
}

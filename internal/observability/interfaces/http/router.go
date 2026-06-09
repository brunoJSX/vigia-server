package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /monitors", h.CreateMonitor)
	mux.HandleFunc("POST /monitors/{id}/pause", h.PauseMonitor)
	mux.HandleFunc("POST /monitors/{id}/resume", h.ResumeMonitor)
	mux.HandleFunc("POST /monitors/{id}/disable", h.DisableMonitor)
	mux.HandleFunc("GET /monitors/{id}/history", h.MonitorHistory)

	return mux
}

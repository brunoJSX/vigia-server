package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /monitors", h.CreateMonitor)
	mux.HandleFunc("GET /monitors", h.ListMonitors)
	mux.HandleFunc("POST /monitors/{id}/pause", h.PauseMonitor)
	mux.HandleFunc("POST /monitors/{id}/resume", h.ResumeMonitor)
	mux.HandleFunc("POST /monitors/{id}/disable", h.DisableMonitor)
	mux.HandleFunc("GET /monitors/history", h.AggregateHistory)
	mux.HandleFunc("GET /monitors/{id}/history", h.MonitorHistory)
	mux.HandleFunc("GET /incidents", h.ListIncidents)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

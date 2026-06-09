package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /account", h.GetAccount)
	mux.HandleFunc("PATCH /account", h.UpdateAccount)
	return mux
}

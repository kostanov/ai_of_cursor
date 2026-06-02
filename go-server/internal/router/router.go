package router

import (
	"net/http"

	"ai-of-cursor/go-server/internal/handler"
)

// New регистрирует маршруты API.
func New(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/adduser", h.AddUser)
	mux.HandleFunc("/user/", h.GetUser)
	mux.HandleFunc("/activate/", h.Activate)
	mux.HandleFunc("/slow", h.Slow)
	mux.HandleFunc("/wrong", h.Wrong)
	return mux
}

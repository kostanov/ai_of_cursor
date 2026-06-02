package handler

import (
	"net/http"

	"ai-of-cursor/go-server/internal/config"
	"ai-of-cursor/go-server/internal/models"
	"ai-of-cursor/go-server/internal/service"
	"ai-of-cursor/go-server/internal/util"
)

// Slow — GET /slow.
func (h *Handler) Slow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.MethodNotAllowed(w, http.MethodGet)
		return
	}

	go service.SlowTask(config.SlowTaskIterations)
	util.WriteJSON(w, http.StatusAccepted, models.StatusResponse{Status: "scheduled"})
}

// Wrong — GET /wrong.
func (h *Handler) Wrong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.MethodNotAllowed(w, http.MethodGet)
		return
	}

	util.WriteJSON(w, http.StatusInternalServerError, models.WrongErrorBody{
		Msg:    "error",
		Detail: "division by zero",
	})
}

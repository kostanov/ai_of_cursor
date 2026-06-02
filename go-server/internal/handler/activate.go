package handler

import (
	"net/http"
	"time"

	"ai-of-cursor/go-server/internal/config"
	"ai-of-cursor/go-server/internal/models"
	"ai-of-cursor/go-server/internal/util"
)

// Activate — GET|POST /activate/{uid}.
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		util.MethodNotAllowed(w, http.MethodGet, http.MethodPost)
		return
	}

	uid, err := util.ParseIDFromPath(r.URL.Path, "/activate/")
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, models.ErrorBody{Error: "invalid_user_id"})
		return
	}

	time.Sleep(config.ActivateDelay)
	active := h.active.Add(uid)

	util.WriteJSON(w, http.StatusOK, models.ActivateResponse{
		Status: "ok",
		Active: active,
	})
}

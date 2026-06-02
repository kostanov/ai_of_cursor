package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-of-cursor/go-server/internal/models"
	"ai-of-cursor/go-server/internal/repository"
	"ai-of-cursor/go-server/internal/util"
)

// AddUser — POST /adduser.
func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.MethodNotAllowed(w, http.MethodPost)
		return
	}

	var req models.AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, models.ErrorBody{Error: "invalid json"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		util.WriteJSON(w, http.StatusBadRequest, models.ErrorBody{Error: "name is required"})
		return
	}

	userID, err := h.users.Create(name)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, models.ErrorBody{Error: "db_insert_failed"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, models.AddUserResponse{
		Status: "ok",
		ID:     userID,
		Name:   name,
	})
}

// GetUser — GET /user/{uid}.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.MethodNotAllowed(w, http.MethodGet)
		return
	}

	uid, err := util.ParseIDFromPath(r.URL.Path, "/user/")
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, models.ErrorBody{Error: "invalid_user_id"})
		return
	}

	user, err := h.users.GetByID(uid)
	if repository.ErrNotFound(err) {
		util.WriteJSON(w, http.StatusNotFound, models.ErrorBody{Error: "not_found"})
		return
	}
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, models.ErrorBody{Error: "db_query_failed"})
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}

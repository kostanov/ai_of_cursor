package handler

import (
	"ai-of-cursor/go-server/internal/repository"
	"ai-of-cursor/go-server/internal/service"
)

// Handler — HTTP-обработчики API.
type Handler struct {
	users  *repository.UserRepository
	active *service.ActiveUsers
}

// New создаёт Handler с зависимостями.
func New(users *repository.UserRepository, active *service.ActiveUsers) *Handler {
	return &Handler{users: users, active: active}
}

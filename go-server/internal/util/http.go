package util

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// WriteJSON отправляет JSON-ответ с заданным статусом.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json: %v", err)
	}
}

// MethodNotAllowed — 405 с заголовком Allow.
func MethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
}

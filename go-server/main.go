package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const (
	dbPath         = "test.db"
	maxActiveUsers = 100
)

type server struct {
	db *sql.DB

	activeMu    sync.Mutex
	activeUsers []int
}

type addUserRequest struct {
	Name string `json:"name"`
}

func main() {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err = initDB(db); err != nil {
		log.Fatalf("init db: %v", err)
	}

	srv := &server{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/adduser", srv.handleAddUser)
	mux.HandleFunc("/user/", srv.handleGetUser)
	mux.HandleFunc("/activate/", srv.handleActivate)
	mux.HandleFunc("/slow", srv.handleSlow)
	mux.HandleFunc("/wrong", srv.handleWrong)

	addr := "0.0.0.0:8080"
	log.Printf("server started on %s", addr)
	if err = http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)",
	)
	return err
}

func (s *server) handleAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req addUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	result, err := s.db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db_insert_failed"})
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db_last_insert_id_failed"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "ok",
		"id":     userID,
		"name":   name,
	})
}

func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	uid, err := parseIDFromPath(r.URL.Path, "/user/")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_user_id"})
		return
	}

	var id int
	var name string
	err = s.db.QueryRow("SELECT id, name FROM users WHERE id = ?", uid).Scan(&id, &name)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db_query_failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"id": id, "name": name})
}

func (s *server) handleActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodGet, http.MethodPost)
		return
	}

	uid, err := parseIDFromPath(r.URL.Path, "/activate/")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_user_id"})
		return
	}

	time.Sleep(100 * time.Millisecond)

	s.activeMu.Lock()
	s.activeUsers = append(s.activeUsers, uid)
	if len(s.activeUsers) > maxActiveUsers {
		overflow := len(s.activeUsers) - maxActiveUsers
		s.activeUsers = s.activeUsers[overflow:]
	}
	activeCopy := append([]int(nil), s.activeUsers...)
	s.activeMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"active": activeCopy,
	})
}

func (s *server) handleSlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	go slowTask(200000)
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "scheduled"})
}

func (s *server) handleWrong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"msg":    "error",
		"detail": "division by zero",
	})
}

func slowTask(iterations int) int {
	x := 0
	for i := 0; i < iterations; i++ {
		x += i * (i + 1) / 2
	}
	return x
}

func parseIDFromPath(path, prefix string) (int, error) {
	rawID := strings.TrimPrefix(path, prefix)
	if rawID == "" || strings.Contains(rawID, "/") {
		return 0, errors.New("invalid id")
	}
	uid, err := strconv.Atoi(rawID)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func methodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json: %v", err)
	}
}

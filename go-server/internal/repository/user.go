package repository

import (
	"database/sql"
	"errors"

	"ai-of-cursor/go-server/internal/models"
)

// UserRepository — доступ к таблице users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт репозиторий.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create добавляет пользователя и возвращает id.
func (r *UserRepository) Create(name string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetByID возвращает пользователя или sql.ErrNoRows.
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, name FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ErrNotFound — обёртка для отсутствующей записи.
func ErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

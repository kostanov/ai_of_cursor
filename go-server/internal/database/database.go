package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open открывает SQLite по пути.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

// InitSchema создаёт таблицу users при отсутствии.
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)",
	)
	return err
}

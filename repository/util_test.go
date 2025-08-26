package repository

import (
	"database/sql"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// getTestRepository sets up an in-memory SQLite database and returns a Repository instance for testing
func getTestRepository() *Repository {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	repo, err := NewRepository(
		db,
		filepath.Join("..", "assets", "queries"), // Use relative path to queries directory
	)
	if err != nil {
		panic(err)
	}
	return repo
}

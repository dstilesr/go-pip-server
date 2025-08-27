package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewRepository creates a new Repository instance
func NewRepository(db *sql.DB, queriesPath string) (*Repository, error) {
	info, err := os.Stat(queriesPath)

	// Verify query directory - must be directory and must contain required files
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("queries path does not exist: %w", err)
		}
		return nil, fmt.Errorf("error accessing queries path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("queries path is not a directory")
	}
	err = checkQueryDirectory(queriesPath)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		DB:          db,
		QueriesPath: queriesPath,
	}
	return repo, nil
}

// checkQueryDirectory verifies that all required SQL files are present in the specified directory
func checkQueryDirectory(path string) error {
	req := strings.Split(RequiredQueryFiles, ";")
	for _, f := range req {
		fullPath := filepath.Join(path, f)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("required query file %s does not exist in %s", f, path)
		} else if err != nil {
			return fmt.Errorf("error accessing file %s: %w", f, err)
		}
	}
	return nil
}

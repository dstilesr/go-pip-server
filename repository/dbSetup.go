package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SetUpDB reads the table creation SQL file and executes it against the database
// to set up the necessary tables. It uses a timeout to prevent long-running operations.
func (r *Repository) SetUpDB() error {
	schemaPath := filepath.Join(r.QueriesPath, "table-creation.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TableCreationTimeoutSeconds*time.Second)
	defer cancel()

	_, err = r.DB.ExecContext(ctx, string(schemaSQL))
	if err != nil {
		return fmt.Errorf("error executing schema SQL: %w", err)
	}
	return nil
}

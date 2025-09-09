package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

// GetOrCreateProject retrieves a project by name, or creates it if it does not exist.
func (r *Repository) GetOrCreateProject(n string, c context.Context) (*Project, error) {
	tx, err := r.DB.BeginTx(c, nil)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"insert into projects (name) values (?) on conflict(name) do nothing",
		n,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var p Project
	err = r.DB.QueryRow("select id, name from projects where name = ?", n).Scan(&p.ID, &p.Name)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetAllProjects retrieves all projects from the database along with the highest project ID.
func (r *Repository) GetAllProjects(c context.Context) (*AllProjects, error) {
	rows, err := r.DB.QueryContext(c, "select id, name from projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projects := make([]*Project, 0, 64)
	var maxId int64 = -1

	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			slog.Error("Failed to scan project row", "error", err)
		}
		if p.ID > maxId {
			maxId = p.ID
		}
		projects = append(projects, &p)
	}
	return &AllProjects{
		Projects:   projects,
		LastSerial: maxId,
	}, nil
}

// GetLatestProjectVersionId retrieves the latest version ID for a given project name,
// optionally within a transaction.
func (r *Repository) GetLatestProjectVersionId(pn string, c context.Context, tx *sql.Tx) (int64, error) {
	qry := `select v.id 
            from versions as v 
            join projects as p on v.project_id = p.id 
            where p.name = ? 
            order by v.created_at desc, v.id desc
            limit 1`
	var id int64
	var err error = nil
	if tx != nil {
		err = tx.QueryRowContext(c, qry, pn).Scan(&id)
	} else {
		err = r.DB.QueryRowContext(c, qry, pn).Scan(&id)
	}
	return id, err
}

// CreateProjectVersion creates a new project version along with its metadata in the database.
func (r *Repository) CreateProjectVersion(pvi *ProjectVersionInsert, c context.Context) error {
	// Get / create parent project
	proj, err := r.GetOrCreateProject(pvi.ProjectName, c)
	if err != nil {
		slog.Error("Unable to get or create project", "error", err)
		return err
	}

	// Insert new version
	tx, err := r.DB.BeginTx(c, nil)
	if err != nil {
		slog.Error("Unable to begin transaction", "error", err)
		return err
	}
	_, err = tx.ExecContext(
		c,
		"insert into versions (project_id, digest, digest_type, filepath, version, file_type) values (?, ?, ?, ?, ?, ?)",
		proj.ID,
		pvi.Digest,
		pvi.DigestType,
		pvi.FilePath,
		pvi.Version,
		pvi.FileType,
	)
	if err != nil {
		slog.Error("Unable to insert version", "error", err)
		tx.Rollback()
		return err
	}
	vId, err := r.GetLatestProjectVersionId(pvi.ProjectName, c, tx)
	if err != nil {
		slog.Error("Unable to get latest version ID", "error", err)
		tx.Rollback()
		return err
	}

	// Add metadata fields
	if len(pvi.Metadata) > 0 {
		metaQry := makeMetaInsertQuery(vId, len(pvi.Metadata))
		flatKVs := flattenKVs(pvi.Metadata)
		_, err = tx.ExecContext(c, metaQry, flatKVs...)
		if err != nil {
			slog.Error("Unable to insert metadata", "error", err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// makeMetaInsertQuery constructs an SQL insert query template for metadata key-value pairs.
func makeMetaInsertQuery(vId int64, total int) string {
	base := "insert into version_metadata_fields (version_id, key, value) values "
	slots := make([]string, 0, total)
	for range total {
		slots = append(slots, fmt.Sprintf("(%d, ?, ?)", vId))
	}
	base += strings.Join(slots, ", ")
	return base
}

// flattenKVs converts a slice of KeyVal structs into a flat slice of "any" for SQL insertion.
func flattenKVs(kvs []*KeyVal) []any {
	out := make([]any, 0, len(kvs)*2)
	for _, kv := range kvs {
		out = append(out, kv.Key, kv.Val)
	}
	return out
}

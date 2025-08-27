package repository

import (
	"context"
	"log/slog"
)

// Project represents a project entity in the database.
type Project struct {
	ID   int64  `json:"_last-serial"`
	Name string `json:"name"`
}

// AllProjects represents a collection of all projects along with the last serial number.
type AllProjects struct {
	LastSerial int64
	Projects   []*Project
}

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

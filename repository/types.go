package repository

import "database/sql"

// Project represents a project entity in the database.
type Project struct {
	ID   int64  `json:"_last-serial"`
	Name string `json:"name"`
}

// AllProjects represents a collection of all projects along with the last serial number.
// This will be used to return the response for the /simple/ endpoint.
type AllProjects struct {
	LastSerial int64
	Projects   []*Project
}

// Repository holds the DB connection pool and performs database operations.
type Repository struct {
	DB          *sql.DB
	QueriesPath string
}

// KeyVal represents a key-value pair, used for metadata storage.
type KeyVal struct {
	Key string
	Val string
}

// ProjectVersionInsert represents a specific version of a project to store in the
// database.
type ProjectVersionInsert struct {
	ProjectName string
	Version     string
	Digest      string
	DigestType  string
	FilePath    string
	FileType    string
	Metadata    []*KeyVal
}

// ProjectFileMeta represents metadata associated with a project file
type ProjectFileMeta struct{}

package main

import (
	"database/sql"
	"flag"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Config is the overall application configuration
type Config struct {
	SQLiteFile    string
	Port          int
	HostAddr      string
	QueriesSource string
	DataPath      string
}

// SetUp Parses command-line flags and returns the application configuration
func SetUp() *Config {
	var cfg Config
	flag.StringVar(
		&cfg.SQLiteFile,
		"sqlite-file",
		filepath.Join("_data", "meta.sqlite"),
		"Path to the SQLite database file",
	)
	flag.IntVar(&cfg.Port, "port", 8080, "Port to run the server on")
	flag.StringVar(&cfg.HostAddr, "host-addr", "0.0.0.0", "Host address to bind the server to")
	flag.StringVar(
		&cfg.QueriesSource,
		"queries-dir",
		filepath.Join("assets", "queries"),
		"Directory containing SQL query files",
	)
	flag.StringVar(
		&cfg.DataPath,
		"data-path",
		filepath.Join("_data", "packages"),
		"Path to the directory containing package files",
	)
	flag.Parse()
	return &cfg
}

func main() {
	cfg := SetUp()
	sqlDb, err := sql.Open("sqlite", cfg.SQLiteFile)
	if err != nil {
		panic(err)
	}
	defer sqlDb.Close()

	// Start server and listen for requests
	server, err := NewPipServer(sqlDb, cfg)
	if err != nil {
		log.Fatalf("Error setting up server: %v", err)
	}
	log.Fatal(server.Serve())
}

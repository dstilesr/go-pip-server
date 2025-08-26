package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Config is the overall application configuration
type Config struct {
	SQLiteFile string
	Port       int
	HostAddr   string
}

func SetUp() *Config {
	var cfg Config
	flag.StringVar(&cfg.SQLiteFile, "sqlite-file", "_data/data.db", "Path to the SQLite database file")
	flag.IntVar(&cfg.Port, "port", 8080, "Port to run the server on")
	flag.StringVar(&cfg.HostAddr, "host-addr", "0.0.0.0", "Host address to bind the server to")
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
	server := NewPipServer(
		fmt.Sprintf("%s:%d", cfg.HostAddr, cfg.Port),
		sqlDb,
	)
	log.Fatal(server.Serve())
}

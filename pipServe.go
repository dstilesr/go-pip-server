package main

import (
	"database/sql"
	"errors"
	"fmt"
	"go-pip-server/repository"
	"log/slog"
	"net/http"
	"os"
)

type PipServer struct {
	Server   *http.Server
	DBConn   *sql.DB
	Repo     *repository.Repository
	isSetUp  bool
	DataPath string
}

// NewPipServer Instantiates and sets up a new Pip Server
func NewPipServer(db *sql.DB, cfg *Config) (*PipServer, error) {
	repo, err := repository.NewRepository(db, cfg.QueriesSource)
	if err != nil {
		return nil, err
	}

	err = repo.SetUpDB()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(cfg.DataPath); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.DataPath, 0755)
		if err != nil {
			slog.Error("Error creating data path directory", "error", err)
			os.Exit(1)
		}
	}

	srv := http.Server{Addr: fmt.Sprintf("%s:%d", cfg.HostAddr, cfg.Port)}
	pip := &PipServer{
		Server:   &srv,
		DBConn:   db,
		isSetUp:  false,
		Repo:     repo,
		DataPath: cfg.DataPath,
	}
	err = pip.SetUpRoutes()
	if err != nil {
		return nil, err
	}

	return pip, nil
}

// SetUpRoutes Configures the HTTP routes for the Pip Server
// this should be called only once during initialization
func (p *PipServer) SetUpRoutes() error {
	if p.isSetUp {
		return errors.New("the server routes have already been set up")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/simple/", p.HandleSimpleIndex)
	p.isSetUp = true
	p.Server.Handler = mux

	slog.Info("Server routes have been set up")
	return nil
}

// Serve Starts the HTTP server and listens for incoming requests
func (p *PipServer) Serve() error {
	if !p.isSetUp {
		return errors.New("the server routes have not been set up")
	}
	return p.Server.ListenAndServe()
}

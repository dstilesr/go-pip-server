package main

import (
	"database/sql"
	"errors"
	"go-pip-server/repository"
	"log/slog"
	"net/http"
)

type PipServer struct {
	Server  *http.Server
	DBConn  *sql.DB
	Repo    *repository.Repository
	isSetUp bool
}

// NewPipServer Instantiates and sets up a new Pip Server
func NewPipServer(addr string, db *sql.DB, qp string) (*PipServer, error) {
	repo, err := repository.NewRepository(db, qp)
	if err != nil {
		return nil, err
	}

	srv := http.Server{Addr: addr}
	pip := &PipServer{
		Server:  &srv,
		DBConn:  db,
		isSetUp: false,
		Repo:    repo,
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

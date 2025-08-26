package main

import (
	"database/sql"
	"errors"
	"net/http"
)

type PipServer struct {
	Server  *http.Server
	DBConn  *sql.DB
	isSetUp bool
}

// NewPipServer Instantiates and sets up a new Pip Server
func NewPipServer(addr string, db *sql.DB) *PipServer {
	srv := http.Server{Addr: addr}
	pip := &PipServer{
		Server:  &srv,
		DBConn:  db,
		isSetUp: false,
	}
	pip.SetUpRoutes()
	return pip
}

// SetUpRoutes Configures the HTTP routes for the Pip Server
// this should be called only once during initialization
func (p *PipServer) SetUpRoutes() error {
	if p.isSetUp {
		return errors.New("the server routes have already been set up")
	}
	// TODO - specify endpoints
	p.isSetUp = true
	return nil
}

// Serve Starts the HTTP server and listens for incoming requests
func (p *PipServer) Serve() error {
	if !p.isSetUp {
		return errors.New("the server routes have not been set up")
	}
	return p.Server.ListenAndServe()
}

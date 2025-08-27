package main

import (
	"encoding/json"
	"go-pip-server/repository"
	"log/slog"
	"net/http"
	"strconv"
)

// APIMeta Returned with Simple API responses
type APIMeta struct {
	Version string `json:"version"`
	MaxId   int64  `json:"_last-serial"`
}

// SimpleIdxResponse Returned by the /simple/ endpoint
type SimpleIdxResponse struct {
	Metadata APIMeta               `json:"meta"`
	Projects []*repository.Project `json:"projects"`
}

func ParseForm(w http.ResponseWriter, r *http.Request) {}

func (p *PipServer) HandleSimpleIndex(w http.ResponseWriter, r *http.Request) {
	ps, err := p.Repo.GetAllProjects(r.Context())
	if err != nil {
		slog.Error("Error fetching projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rsp := SimpleIdxResponse{
		Projects: ps.Projects,
		Metadata: APIMeta{
			Version: APIVersion,
			MaxId:   ps.LastSerial,
		},
	}
	w.Header().Set("Content-Type", "application/vnd.pypi.simple.v1+json")
	w.Header().Set("X-PyPI-Last-Serial", strconv.FormatInt(ps.LastSerial, 10))
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&rsp)
	if err != nil {
		slog.Error("Error encoding JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

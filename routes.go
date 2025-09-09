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

// HandleUpload parses multipart form data from an HTTP request to upload a package.
func (p *PipServer) HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, `{"detail": "Invalid form data"}`, http.StatusBadRequest)
		return
	}

	pvi, err := p.PrepareFormData(r.MultipartForm)
	if err != nil {
		slog.Error("Error preparing form data", "error", err)
		http.Error(w, `{"detail": "Invalid form data"}`, http.StatusBadRequest)
		return
	}

	err = p.Repo.CreateProjectVersion(pvi, r.Context())
	if err != nil {
		slog.Error("Error inserting project version", "error", err)
		http.Error(w, `{"detail": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (p *PipServer) HandleSimpleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != JSONHeader {
		http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
	}

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
	w.Header().Set("Content-Type", JSONHeader)
	w.Header().Set("X-PyPI-Last-Serial", strconv.FormatInt(ps.LastSerial, 10))
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&rsp)
	if err != nil {
		slog.Error("Error encoding JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

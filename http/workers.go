package http

import (
	"net/http"

	"github.com/travis-ci/junction/junction"
)

func handleWorkers(core *junction.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleWorkersPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleWorkersPost(core *junction.Core, w http.ResponseWriter, r *http.Request) {
	// Get authentication token
	token := requestAuth(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, nil)
		return
	}

	// Parse the request
	var req WorkerCreateRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	id, err := core.WorkerHandler.Create(token, req.Queue, req.MaxJobCount, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, &WorkerCreateResponse{ID: id})
}

type WorkerCreateRequest struct {
	Queue       string `json:"queue"`
	MaxJobCount int    `json:"max-job-count"`
}

type WorkerCreateResponse struct {
	ID string `json:"id"`
}

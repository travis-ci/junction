package http

import (
	"net/http"

	"github.com/travis-ci/junction/junction"
)

func handleWorkers(core *junction.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleWorkersGet(core, w, r)
		case "POST":
			handleWorkersPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleWorkersGet(core *junction.Core, w http.ResponseWriter, r *http.Request) {
	// Get authentication token
	token := requestAuth(r)
	if token == "" {
		respondError(w, http.StatusUnauthorized, nil)
		return
	}

	workers, err := core.AdminHandler.ListWorkers(token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	resp := WorkerListResponse{Workers: []WorkerListResponseWorker{}}
	for _, worker := range workers {
		resp.Workers = append(resp.Workers, WorkerListResponseWorker{
			ID:          worker.ID,
			Queue:       worker.Queue,
			MaxJobCount: worker.MaxJobCount,
			Attributes:  worker.Attributes,
		})
	}

	respondOk(w, resp)
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

	id, err := core.WorkerHandler.Create(token, req.Queue, req.MaxJobCount, req.Attributes)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, &WorkerCreateResponse{ID: id})
}

type WorkerCreateRequest struct {
	Queue       string            `json:"queue"`
	MaxJobCount int               `json:"max-job-count"`
	Attributes  map[string]string `json:"attributes"`
}

type WorkerListResponseWorker struct {
	ID          string            `json:"id"`
	Queue       string            `json:"queue"`
	MaxJobCount int               `json:"max-job-count"`
	Attributes  map[string]string `json:"attributes"`
}

type WorkerListResponse struct {
	Workers []WorkerListResponseWorker `json:"workers"`
}

type WorkerCreateResponse struct {
	ID string `json:"id"`
}

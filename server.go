package junction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/braintree/manners"
	"github.com/gorilla/mux"
)

type server struct {
	addr string
	r    *mux.Router

	workerRepo WorkerRepository
}

func newServer(addr string, workerRepo WorkerRepository) *server {
	return &server{
		addr:       addr,
		r:          mux.NewRouter(),
		workerRepo: workerRepo,
	}
}

func (srv *server) Setup() {
	srv.setupRoutes()
}

func (srv *server) Run() {
	manners.ListenAndServe(srv.addr, srv.r)
}

func (srv *server) setupRoutes() {
	srv.r.HandleFunc(`/workers`, srv.handleWorkersPost).Methods("POST").Name("workers-post")
	srv.r.HandleFunc(`/workers/{id}/heartbeat`, srv.handleWorkersHeartbeat).Methods("POST").Name("workers-heartbeat")
}

func (srv *server) handleWorkersPost(w http.ResponseWriter, req *http.Request) {
	var parsedRequest struct {
		Queue       string `json:"queue"`
		MaxJobCount int    `json:"max-job-count"`
	}

	err := json.NewDecoder(req.Body).Decode(&parsedRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"@type":"error","error_message":"Invalid JSON"}`)
		return
	}

	if parsedRequest.Queue == "" || parsedRequest.MaxJobCount < 1 {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"@type":"error","error_message":"queue must not be empty and max-job-count must be >0"}`)
		return
	}

	worker := Worker{
		ID:          uuid.NewRandom(),
		Queue:       parsedRequest.Queue,
		MaxJobCount: parsedRequest.MaxJobCount,
	}

	err = srv.workerRepo.Create(worker)
	if err != nil {
		// TODO: Send to Sentry
		log.Printf("error storing worker: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"@type":"error","error_message":"couldn't create worker"}`)
		return
	}

	w.WriteHeader(http.StatusCreated)

	var responseWorker struct {
		ID                 string   `json:"id"`
		Queue              string   `json:"queue"`
		MaxJobCount        int      `json:"max-job-count"`
		CurrentAssignments []string `json:"current-assignments"`
	}
	responseWorker.ID = worker.ID.String()
	responseWorker.Queue = worker.Queue
	responseWorker.MaxJobCount = worker.MaxJobCount

	err = json.NewEncoder(w).Encode(responseWorker)
	if err != nil {
		// TODO: Log? Send to Sentry?
	}
}

func (srv *server) handleWorkersHeartbeat(w http.ResponseWriter, req *http.Request) {
	var parsedRequest struct {
		CurrentAssignments []string `json:"current-assignments"`
	}

	err := json.NewDecoder(req.Body).Decode(&parsedRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"@type":"error","error_message":"Invalid JSON"}`)
		return
	}

	vars := mux.Vars(req)
	id := vars["id"]

	worker, err := srv.workerRepo.Fetch(uuid.Parse(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"@type":"error","error_message":"No worker with that ID found"}`)
		return
	}

	now := time.Now()
	worker.LastHeartbeat = &now

	err = srv.workerRepo.Update(worker)
	if err != nil {
		log.Printf("error storing worker: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"@type":"error","error_message":"couldn't process heartbeat"}`)
		return
	}

}

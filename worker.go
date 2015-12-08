package junction

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Worker struct {
	ID            uuid.UUID
	Queue         string
	LastHeartbeat *time.Time
	MaxJobCount   int
}

// WorkerHeartbeatHandler handles all the business logic behind handling a
// heartbeat from a worker.
type WorkerHeartbeatHandler struct {
	repo WorkerRepository
}

// Heartbeat handles the heartbeat from a given worker. It takes the worker,
// the assignments the worker claims to be working on and returns a list of
// assignments for the worker to work on, or an error.
func (whh *WorkerHeartbeatHandler) Heartbeat(w Worker, assignments []string) ([]string, error) {
	now := time.Now()
	w.LastHeartbeat = &now
	err := whh.repo.Store(w)
	return []string{}, err
}

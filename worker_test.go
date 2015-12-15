package junction

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"code.google.com/p/go-uuid/uuid"
)

func TestWorkerHeartbeat(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	whh := &WorkerHeartbeatHandler{repo: mwr}
	worker := Worker{
		ID:            uuid.NewRandom(),
		Queue:         "test-queue",
		LastHeartbeat: nil,
		MaxJobCount:   10,
	}
	mwr.Create(worker)

	assignments, err := whh.Heartbeat(worker, []string{})
	assert.Nil(t, err)
	assert.NotNil(t, assignments, "assignments should be a slice (empty or not) if err is nil")

	storedWorker, _ := mwr.Fetch(worker.ID)
	assert.NotNil(t, storedWorker.LastHeartbeat, "LastHeartbeat should be updated")
}

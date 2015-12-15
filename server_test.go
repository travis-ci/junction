package junction

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"code.google.com/p/go-uuid/uuid"
)

func TestServerWorkerPost(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	body := bytes.NewReader([]byte(`{"queue":"test-queue","max-job-count":10}`))
	req, err := http.NewRequest("POST", ts.URL+"/workers", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	var responseWorker struct {
		ID                 string   `json:"id"`
		Queue              string   `json:"queue"`
		MaxJobCount        int      `json:"max-job-count"`
		CurrentAssignments []string `json:"current-assignments"`
	}

	err = json.NewDecoder(res.Body).Decode(&responseWorker)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusCreated, "unexpected status code")
	assert.Equal(t, responseWorker.Queue, "test-queue", "queue in reply is incorrect")
	assert.Equal(t, responseWorker.MaxJobCount, 10, "max job count in reply is incorrect")
	assert.Empty(t, responseWorker.CurrentAssignments, "worker shouldn't have assignments initially")

	storedWorker, ok := mwr.workers[responseWorker.ID]
	assert.True(t, ok, "worker wasn't stored in repository")
	assert.Equal(t, storedWorker.Queue, "test-queue", "queue in database is incorrect")
	assert.Equal(t, storedWorker.MaxJobCount, 10, "max job count in database is incorrect")
}

func TestServerWorkerPostInvalidRequest(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	body := bytes.NewReader([]byte(`{thisisnotvalidjson`))
	req, err := http.NewRequest("POST", ts.URL+"/workers", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusBadRequest, "unexpected status code")
	assert.Equal(t, len(mwr.workers), 0, "no workers should be stored")
}

func TestServerWorkerPostMissingQueue(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	body := bytes.NewReader([]byte(`{"max-job-count":10}`))
	req, err := http.NewRequest("POST", ts.URL+"/workers", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, 422, "unexpected status code")
}

func TestServerWorkerPostMaxJobCountTooLow(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	body := bytes.NewReader([]byte(`{"max-job-count":0,"queue":"test-queue"}`))
	req, err := http.NewRequest("POST", ts.URL+"/workers", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, 422, "unexpected status code")
}

func TestServerWorkerHeartbeat(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	worker := Worker{
		ID:            uuid.NewRandom(),
		Queue:         "test-queue",
		LastHeartbeat: nil,
		MaxJobCount:   10,
	}
	mwr.Create(worker)

	body := bytes.NewReader([]byte(`{"current-assignments":[]}`))
	req, err := http.NewRequest("POST", ts.URL+"/workers/"+worker.ID.String()+"/heartbeat", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK, "unexpected status code")

	storedWorker, ok := mwr.workers[worker.ID.String()]
	assert.True(t, ok, "worker disappeared from database")
	assert.NotNil(t, storedWorker.LastHeartbeat, "heartbeat should be updated")
}

func TestServerWorkerHeartbeatNotFound(t *testing.T) {
	mwr := &MapWorkerRepository{workers: make(map[string]Worker)}
	srv := newServer("", mwr)
	srv.Setup()
	ts := httptest.NewServer(srv.r)
	defer ts.Close()

	body := bytes.NewReader([]byte(`{"current-assignments":[]}`))
	req, err := http.NewRequest("POST", ts.URL+"/workers/bogus-id/heartbeat", body)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusNotFound, "unexpected status code")
}

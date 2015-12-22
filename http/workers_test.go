package http

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travis-ci/junction/junction"
)

func TestWorkers_get(t *testing.T) {
	core := junction.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	id, err := core.WorkerHandler.Create("worker-token-1", "test-queue", 10, nil)
	require.NoError(t, err)

	resp := testHttpGet(t, "admin-token-1", addr+"/workers")

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	require.Contains(t, actual, "workers")

	require.Equal(t, []interface{}{
		map[string]interface{}{
			"id":            id,
			"queue":         "test-queue",
			"max-job-count": float64(10),
			"attributes":    nil,
		},
	}, actual["workers"])
}

func TestWorkers_post(t *testing.T) {
	core := junction.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPost(t, "worker-token-1", addr+"/workers", map[string]interface{}{
		"queue":         "test-queue",
		"max-job-count": 10,
		"attributes": map[string]string{
			"version": "1.0.0",
		},
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	require.Contains(t, actual, "id")
}

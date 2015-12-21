package http

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travis-ci/junction/junction"
)

func TestWorkers_post(t *testing.T) {
	core := junction.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPost(t, "worker-token-1", addr+"/workers", map[string]interface{}{
		"queue":         "test-queue",
		"max-job-count": 10,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	require.Contains(t, actual, "id")
}

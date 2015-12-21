package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
)

func testHttpPost(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "POST", token, addr, body)
}

func testHttpData(t *testing.T, method string, token string, addr string, body interface{}) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		err := enc.Encode(body)
		require.Nil(t, err)
	}

	req, err := http.NewRequest(method, addr, bodyReader)
	require.Nil(t, err)

	req.Header.Set("Content-Type", "application/json")

	if len(token) != 0 {
		req.Header.Set("X-Junction-Token", token)
	}

	client := cleanhttp.DefaultClient()
	client.Timeout = 60 * time.Second

	// From https://github.com/michiwend/gomusicbrainz/pull/4/files
	defaultRedirectLimit := 30

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > defaultRedirectLimit {
			return fmt.Errorf("%d consecutive requests(redirects)", len(via))
		}
		if len(via) == 0 {
			// No redirects
			return nil
		}
		return nil
	}

	resp, err := client.Do(req)
	require.Nil(t, err)

	return resp
}

func testResponseStatus(t *testing.T, resp *http.Response, code int) {
	if resp.StatusCode != code {
		body := new(bytes.Buffer)
		io.Copy(body, resp.Body)
		resp.Body.Close()

		t.Fatalf(
			"Expected status %d, got %d. Body:\n\n%s",
			code,
			resp.StatusCode,
			body.String(),
		)
	}
}

func testResponseBody(t *testing.T, resp *http.Response, out interface{}) {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(out)
	require.Nil(t, err)
}

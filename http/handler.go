package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/travis-ci/junction/junction"
)

// AuthHeaderName is the name of the header containing the token.
const AuthHeaderName = "X-Junction-Token"

// Handler returns an http.Handlre for the API. This can be used on its own to
// mount the Junction API within another web server.
func Handler(core *junction.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/workers", handleWorkers(core))

	return mux
}

func parseRequest(r *http.Request, out interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(out)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Failed to parse JSON input: %s", err)
	}
	return err
}

func requestAuth(r *http.Request) string {
	return r.Header.Get(AuthHeaderName)
}

func respondError(w http.ResponseWriter, status int, err error) {
	if err == junction.ErrAuthenticationError {
		status = http.StatusUnauthorized
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

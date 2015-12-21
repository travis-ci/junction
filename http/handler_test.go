package http

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_error(t *testing.T) {
	w := httptest.NewRecorder()

	respondError(w, 500, errors.New("Test Error"))

	require.Equal(t, 500, w.Code)

}

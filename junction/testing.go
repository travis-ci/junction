package junction

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travis-ci/junction/database"
)

// TestCore returns a pure in-memory core for testing.
func TestCore(t *testing.T) *Core {
	core, err := NewCore(&CoreConfig{
		Database:     database.NewInMem(),
		WorkerTokens: []string{"worker-token-1"},
	})
	require.Nil(t, err)

	return core
}

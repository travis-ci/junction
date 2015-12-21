package junction

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travis-ci/junction/database"
)

func TestWorkerHandlerCreate(t *testing.T) {
	core := TestCore(t)

	id, err := core.WorkerHandler.Create("worker-token-1", "test-queue", 10, nil)
	require.Nil(t, err)

	worker, err := core.database.Get(id)
	require.Nil(t, err)
	require.Equal(t, database.Worker{
		ID:          id,
		Queue:       "test-queue",
		MaxJobCount: 10,
		Attributes:  nil,
	}, worker)
}

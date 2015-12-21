package junction

import (
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
	"github.com/travis-ci/junction/database"
)

func TestAdminHandlerListWorkers(t *testing.T) {
	core := TestCore(t)

	_, err := core.AdminHandler.ListWorkers("invalid-token")
	require.Error(t, err)
	require.Equal(t, err, ErrAuthenticationError)

	workers, err := core.AdminHandler.ListWorkers("admin-token-1")
	require.NoError(t, err)
	require.Empty(t, workers)

	worker := database.Worker{
		ID:          uuid.GenerateUUID(),
		Queue:       "test-queue",
		MaxJobCount: 10,
	}

	err = core.database.CreateWorker(worker)
	require.NoError(t, err)

	workers, err = core.AdminHandler.ListWorkers("admin-token-1")
	require.NoError(t, err)
	require.Len(t, workers, 1)
	require.Contains(t, workers, worker)
}

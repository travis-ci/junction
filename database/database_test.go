package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDatabase(t *testing.T, db Database) {
	workerID := "12345678-1234-4321-abcd-0123456789ab"

	// Worker shouldn't exist to begin with
	worker, err := db.Get(workerID)
	require.Zero(t, worker)
	require.NotNil(t, err)

	// Create a new worker
	worker = Worker{
		ID:          workerID,
		Queue:       "test-queue",
		MaxJobCount: 10,
		Attributes:  nil,
	}
	err = db.Create(worker)
	require.Nil(t, err)

	// Attempt to create a new worker with the same ID
	err = db.Create(worker)
	require.NotNil(t, err)

	// Retrieve the stored worker
	fetchedWorker, err := db.Get(worker.ID)
	require.Nil(t, err)
	require.Equal(t, worker, fetchedWorker)

	// Update the stored worker
	worker.Queue = "new-queue"
	err = db.Update(worker)
	require.Nil(t, err)

	// Check that attributes were updated
	fetchedWorker, err = db.Get(workerID)
	require.Nil(t, err)
	require.Equal(t, worker, fetchedWorker)

	// Delete a worker
	err = db.Delete(worker.ID)
	require.Nil(t, err)

	// Check that we can no longer fetch worker
	fetchedWorker, err = db.Get(workerID)
	require.NotNil(t, err)
	require.Zero(t, fetchedWorker)
}

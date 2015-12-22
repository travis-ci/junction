package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testDatabase(t *testing.T, db Database) {
	workerID := "12345678-1234-4321-abcd-0123456789ab"

	// Worker shouldn't exist to begin with
	worker, err := db.GetWorker(workerID)
	require.Zero(t, worker)
	require.NotNil(t, err)

	// Create a new worker
	worker = Worker{
		ID:          workerID,
		Queue:       "test-queue",
		MaxJobCount: 10,
		Attributes:  map[string]string{"version": "1.0.0"},
	}
	err = db.CreateWorker(worker)
	require.Nil(t, err)

	// Attempt to create a new worker with the same ID
	err = db.CreateWorker(worker)
	require.NotNil(t, err)

	// Create a new worker with a different ID
	worker2 := Worker{
		ID:          "87654321-1234-4321-abcd-0123456789ab",
		Queue:       "test-queue",
		MaxJobCount: 10,
		Attributes:  nil,
	}
	err = db.CreateWorker(worker2)
	require.Nil(t, err)

	// Retrieve the stored worker
	fetchedWorker, err := db.GetWorker(worker.ID)
	require.Nil(t, err)
	require.Equal(t, worker, fetchedWorker)

	// Update the stored worker
	worker.Queue = "new-queue"
	err = db.UpdateWorker(worker)
	require.Nil(t, err)

	// Check that attributes were updated
	fetchedWorker, err = db.GetWorker(workerID)
	require.Nil(t, err)
	require.Equal(t, worker, fetchedWorker)

	// Both created workers is returned in list of workers
	workers, err := db.ListWorkers()
	require.Nil(t, err)
	var workerIDs []string
	for _, worker := range workers {
		workerIDs = append(workerIDs, worker.ID)
	}
	require.Len(t, workerIDs, 2)
	require.Contains(t, workerIDs, worker.ID)
	require.Contains(t, workerIDs, worker2.ID)

	// Delete a worker
	err = db.DeleteWorker(worker.ID)
	require.Nil(t, err)

	// Check that we can no longer fetch worker
	fetchedWorker, err = db.GetWorker(workerID)
	require.NotNil(t, err)
	require.Zero(t, fetchedWorker)
}

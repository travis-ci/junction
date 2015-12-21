package database

import "fmt"

// InMem is an implementation of Database that stores everything in memory.
// Meant to be used for tests that require a database.
type InMem struct {
	workers map[string]Worker
}

// NewInMem returns a new empty InMem database.
func NewInMem() *InMem {
	return &InMem{workers: make(map[string]Worker)}
}

// ListWorkers is used to list all the workers in the database
func (db *InMem) ListWorkers() ([]Worker, error) {
	var workers []Worker

	for _, worker := range db.workers {
		workers = append(workers, worker)
	}

	return workers, nil
}

// CreateWorker is used to store a new worker in the database.
func (db *InMem) CreateWorker(worker Worker) error {
	_, ok := db.workers[worker.ID]
	if ok {
		return fmt.Errorf("worker with ID %s already exists", worker.ID)
	}

	db.workers[worker.ID] = worker

	return nil
}

// GetWorker is used to retrieve a worker that was previously stored in the database.
func (db *InMem) GetWorker(workerID string) (Worker, error) {
	worker, ok := db.workers[workerID]
	if !ok {
		return Worker{}, fmt.Errorf("no worker with ID %s", workerID)
	}

	return worker, nil
}

// UpdateWorker is used to update a previously stored worker in the database.
func (db *InMem) UpdateWorker(worker Worker) error {
	_, ok := db.workers[worker.ID]
	if !ok {
		return fmt.Errorf("no worker with ID %s", worker.ID)
	}

	db.workers[worker.ID] = worker

	return nil
}

// DeleteWorker is used to remove a worker from the database.
func (db *InMem) DeleteWorker(workerID string) error {
	_, ok := db.workers[workerID]
	if !ok {
		return fmt.Errorf("no worker with ID %s", workerID)
	}

	delete(db.workers, workerID)

	return nil
}

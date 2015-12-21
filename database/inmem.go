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

// List is used to list all the workers in the database
func (db *InMem) List() ([]Worker, error) {
	var workers []Worker

	for _, worker := range db.workers {
		workers = append(workers, worker)
	}

	return workers, nil
}

// Create is used to store a new worker in the database.
func (db *InMem) Create(worker Worker) error {
	_, ok := db.workers[worker.ID]
	if ok {
		return fmt.Errorf("worker with ID %s already exists", worker.ID)
	}

	db.workers[worker.ID] = worker

	return nil
}

// Get is used to retrieve a worker that was previously stored in the database.
func (db *InMem) Get(workerID string) (Worker, error) {
	worker, ok := db.workers[workerID]
	if !ok {
		return Worker{}, fmt.Errorf("no worker with ID %s", workerID)
	}

	return worker, nil
}

// Update is used to update a previously stored worker in the database.
func (db *InMem) Update(worker Worker) error {
	_, ok := db.workers[worker.ID]
	if !ok {
		return fmt.Errorf("no worker with ID %s", worker.ID)
	}

	db.workers[worker.ID] = worker

	return nil
}

// Delete is used to remove a worker from the database.
func (db *InMem) Delete(workerID string) error {
	_, ok := db.workers[workerID]
	if !ok {
		return fmt.Errorf("no worker with ID %s", workerID)
	}

	delete(db.workers, workerID)

	return nil
}

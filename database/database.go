package database

// A Database allows persisting workers and fetching them back.
type Database interface {
	ListWorkers() ([]Worker, error)
	CreateWorker(worker Worker) error
	GetWorker(workerID string) (Worker, error)
	UpdateWorker(worker Worker) error
	DeleteWorker(workerID string) error
}

type Worker struct {
	ID          string
	Queue       string
	MaxJobCount int
	Attributes  map[string]string
}

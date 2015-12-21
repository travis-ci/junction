package database

// A Database allows persisting workers and fetching them back.
type Database interface {
	Create(worker Worker) error
	Get(workerID string) (Worker, error)
	Update(worker Worker) error
	Delete(workerID string) error
}

type Worker struct {
	ID          string
	Queue       string
	MaxJobCount int
	Attributes  map[string]string
}

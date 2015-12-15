package junction

import (
	"database/sql"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/lib/pq"
)

// A WorkerRepository can store, fetch and delete Workers
type WorkerRepository interface {
	Fetch(id uuid.UUID) (Worker, error)
	Create(worker Worker) error
	Update(worker Worker) error
	Delete(id uuid.UUID) error
}

// A PostgresWorkerRepository is a WorkerRepository backed by a PostgreSQL
// database.
type PostgresWorkerRepository struct {
	db *sql.DB
}

// A MapWorkerRepository is a WorkerRepository backed by an in-memory map.
// Mostly useful for tests.
type MapWorkerRepository struct {
	workers map[string]Worker
}

func (pwr *PostgresWorkerRepository) Fetch(id uuid.UUID) (Worker, error) {
	worker := Worker{ID: id}

	var lastHeartbeat pq.NullTime
	err := pwr.db.QueryRow(
		"SELECT queue, last_heartbeat, max_job_count FROM junction.workers WHERE id = $1",
		id.String(),
	).Scan(
		&worker.Queue,
		&lastHeartbeat,
		&worker.MaxJobCount,
	)
	if err != nil {
		return Worker{}, err
	}

	if lastHeartbeat.Valid {
		worker.LastHeartbeat = &lastHeartbeat.Time
	}

	return worker, nil
}

func (pwr *PostgresWorkerRepository) Create(worker Worker) error {
	var lastHeartbeat pq.NullTime
	if worker.LastHeartbeat != nil {
		lastHeartbeat.Valid = true
		lastHeartbeat.Time = *worker.LastHeartbeat
	}

	_, err := pwr.db.Exec(
		"INSERT INTO junction.workers (id, queue, created_at, max_job_count, last_heartbeat) VALUES ($1, $2, NOW(), $3, $4)",
		worker.ID.String(),
		worker.Queue,
		worker.MaxJobCount,
		lastHeartbeat,
	)
	return err
}

func (pwr *PostgresWorkerRepository) Update(worker Worker) error {
	var lastHeartbeat pq.NullTime
	if worker.LastHeartbeat != nil {
		lastHeartbeat.Valid = true
		lastHeartbeat.Time = *worker.LastHeartbeat
	}

	_, err := pwr.db.Exec("UPDATE junction.workers SET queue = $1, max_job_count = $2, last_heartbeat = $3 WHERE id = $4", worker.Queue, worker.MaxJobCount, lastHeartbeat, worker.ID.String())
	return err
}

func (pwr *PostgresWorkerRepository) Delete(id uuid.UUID) error {
	_, err := pwr.db.Exec("DELETE FROM junction.workers WHERE id = $1", id.String())
	return err
}

func (mwr *MapWorkerRepository) Fetch(id uuid.UUID) (Worker, error) {
	worker, ok := mwr.workers[id.String()]
	if !ok {
		return Worker{}, fmt.Errorf("no such worker: %v", id)
	}

	return worker, nil
}

func (mwr *MapWorkerRepository) Create(worker Worker) error {
	_, ok := mwr.workers[worker.ID.String()]
	if ok {
		return fmt.Errorf("worker already exists: %v", worker.ID)
	}

	mwr.workers[worker.ID.String()] = worker
	return nil
}

func (mwr *MapWorkerRepository) Update(worker Worker) error {
	_, ok := mwr.workers[worker.ID.String()]
	if !ok {
		return fmt.Errorf("worker doesn't exists: %v", worker.ID)
	}

	mwr.workers[worker.ID.String()] = worker
	return nil
}

func (mwr *MapWorkerRepository) Delete(id uuid.UUID) error {
	delete(mwr.workers, id.String())
	return nil
}

package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

// Postgres is a Database backed by PostgreSQL
type Postgres struct {
	db *sql.DB
}

type PostgresConfig struct {
	URL          string
	MaxOpenConns int
}

// NewPostgres connects to a database given a database URL, and creates a new
// Postgres instance backed by that database.
func NewPostgres(config *PostgresConfig) (*Postgres, error) {
	db, err := sql.Open("postgres", config.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)

	return &Postgres{db: db}, nil
}

// ListWorkers is used to list all the workers in the database
func (db *Postgres) ListWorkers() ([]Worker, error) {
	var workers []Worker

	rows, err := db.db.Query(`SELECT id, queue, max_job_count, attributes FROM junction.workers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var worker Worker
		var attributes hstore.Hstore
		err := rows.Scan(&worker.ID, &worker.Queue, &worker.MaxJobCount, &attributes)
		if err != nil {
			return nil, err
		}

		worker.Attributes = make(map[string]string)
		if attributes.Map != nil {
			for key, value := range attributes.Map {
				if value.Valid {
					worker.Attributes[key] = value.String
				}
			}
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

// CreateWorker is used to store a new worker in the database.
func (db *Postgres) CreateWorker(worker Worker) error {
	var attributes hstore.Hstore
	attributes.Map = make(map[string]sql.NullString)
	if worker.Attributes != nil {
		for key, value := range worker.Attributes {
			attributes.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	_, err := db.db.Exec(
		`INSERT INTO junction.workers (id, queue, max_job_count, attributes, created_at) VALUES ($1, $2, $3, $4, NOW())`,
		worker.ID,
		worker.Queue,
		worker.MaxJobCount,
		attributes,
	)

	return err
}

// GetWorker is used to retrieve a worker that was previously stored in the database.
func (db *Postgres) GetWorker(workerID string) (Worker, error) {
	worker := Worker{ID: workerID}

	var attributes hstore.Hstore

	err := db.db.QueryRow(
		`SELECT queue, max_job_count, attributes FROM junction.workers WHERE id = $1`,
		workerID,
	).Scan(
		&worker.Queue,
		&worker.MaxJobCount,
		&attributes,
	)
	if err != nil {
		return Worker{}, err
	}

	worker.Attributes = make(map[string]string)
	if attributes.Map != nil {
		for key, value := range attributes.Map {
			if value.Valid {
				worker.Attributes[key] = value.String
			}
		}
	}

	return worker, nil
}

// UpdateWorker is used to update a previously stored worker in the database.
func (db *Postgres) UpdateWorker(worker Worker) error {
	var attributes hstore.Hstore
	attributes.Map = make(map[string]sql.NullString)
	if worker.Attributes != nil {
		for key, value := range worker.Attributes {
			attributes.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	_, err := db.db.Exec(
		`UPDATE junction.workers SET queue = $2, max_job_count = $3, attributes = $4 WHERE id = $1`,
		worker.ID,
		worker.Queue,
		worker.MaxJobCount,
		attributes,
	)

	return err
}

// DeleteWorker is used to remove a worker from the database.
func (db *Postgres) DeleteWorker(workerID string) error {
	_, err := db.db.Exec(
		`DELETE FROM junction.workers WHERE id = $1`,
		workerID,
	)

	return err
}

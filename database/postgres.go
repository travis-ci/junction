package database

import (
	"database/sql"

	_ "github.com/lib/pq"
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

// List is used to list all the workers in the database
func (db *Postgres) List() ([]Worker, error) {
	var workers []Worker

	rows, err := db.db.Query(`SELECT id, queue, max_job_count FROM junction.workers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var worker Worker
		err := rows.Scan(&worker.ID, &worker.Queue, &worker.MaxJobCount)
		if err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}

	return workers, nil
}

// Create is used to store a new worker in the database.
func (db *Postgres) Create(worker Worker) error {
	_, err := db.db.Exec(
		`INSERT INTO junction.workers (id, queue, max_job_count, created_at) VALUES ($1, $2, $3, NOW())`,
		worker.ID,
		worker.Queue,
		worker.MaxJobCount,
	)

	return err
}

// Get is used to retrieve a worker that was previously stored in the database.
func (db *Postgres) Get(workerID string) (Worker, error) {
	worker := Worker{ID: workerID}

	err := db.db.QueryRow(
		`SELECT queue, max_job_count FROM junction.workers WHERE id = $1`,
		workerID,
	).Scan(
		&worker.Queue,
		&worker.MaxJobCount,
	)
	if err != nil {
		return Worker{}, err
	}

	return worker, nil
}

// Update is used to update a previously stored worker in the database.
func (db *Postgres) Update(worker Worker) error {
	_, err := db.db.Exec(
		`UPDATE junction.workers SET queue = $2, max_job_count = $3 WHERE id = $1`,
		worker.ID,
		worker.Queue,
		worker.MaxJobCount,
	)

	return err
}

// Delete is used to remove a worker from the database.
func (db *Postgres) Delete(workerID string) error {
	_, err := db.db.Exec(
		`DELETE FROM junction.workers WHERE id = $1`,
		workerID,
	)

	return err
}

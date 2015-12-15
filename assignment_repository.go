package junction

import (
	"database/sql"

	"code.google.com/p/go-uuid/uuid"
)

type Assignment struct {
	ID         uuid.UUID
	WorkerID   uuid.UUID
	JobID      int64
	Dispatched bool
}

// An AssignmentRepository can store, fetch and delete Assignments
type AssignmentRepository interface {
	Fetch(id uuid.UUID) (Assignment, error)
	Create(assignment Assignment) error
	Update(assignment Assignment) error
	Delete(id uuid.UUID) error
}

// A PostgresAssignmentRepository is an AssignmentRepository backed by a
// PostgreSQL database.
type PostgresAssignmentRepository struct {
	db *sql.DB
}

func (par *PostgresAssignmentRepository) Fetch(id uuid.UUID) (Assignment, error) {
	assignment := Assignment{ID: id}

	var workerID string
	err := par.db.QueryRow(
		"SELECT worker_id, job_id, dispatched FROM junction.assignments WHERE id = $1",
		id.String(),
	).Scan(
		&workerID,
		&assignment.JobID,
		&assignment.Dispatched,
	)
	if err != nil {
		return Assignment{}, err
	}

	assignment.WorkerID = uuid.Parse(workerID)

	return assignment, nil
}

func (par *PostgresAssignmentRepository) Create(assignment Assignment) error {
	_, err := par.db.Exec(
		"INSERT INTO junction.assignments (id, worker_id, job_id, dispatched) VALUES ($1, $2, $3, $4)",
		assignment.ID.String(),
		assignment.WorkerID.String(),
		assignment.JobID,
		assignment.Dispatched,
	)
	return err
}

func (par *PostgresAssignmentRepository) Update(assignment Assignment) error {
	_, err := par.db.Exec(
		"UPDATE junction.assignments SET worker_id = $2, job_id = $3, dispatched = $4 WHERE id = $1",
		assignment.ID.String(),
		assignment.WorkerID.String(),
		assignment.JobID,
		assignment.Dispatched,
	)
	return err
}

func (par *PostgresAssignmentRepository) Delete(id uuid.UUID) error {
	_, err := par.db.Exec("DELETE FROM junction.assignments WHERE id = $1", id.String())
	return err
}

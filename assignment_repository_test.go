package junction

import (
	"database/sql"
	"os"
	"testing"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
)

func TestPostgresAssignmentRepository(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer db.Close()
	defer db.Exec("DELETE FROM junction.assignments")
	defer db.Exec("DELETE FROM junction.workers")
	pwr := &PostgresWorkerRepository{db: db}
	par := &PostgresAssignmentRepository{db: db}

	runAssignmentRepositoryTests(t, par, pwr)
}

func runAssignmentRepositoryTests(t *testing.T, repo AssignmentRepository, workerRepo WorkerRepository) {
	Convey("Given a worker that has been stored", t, func() {
		worker := Worker{
			ID:            uuid.NewRandom(),
			Queue:         "test-queue",
			LastHeartbeat: nil,
			MaxJobCount:   10,
		}
		err := workerRepo.Create(worker)
		So(err, ShouldBeNil)

		Convey("Given an assignment that has been created", func() {
			storedAssignment := Assignment{
				ID:         uuid.NewRandom(),
				WorkerID:   worker.ID,
				JobID:      123,
				Dispatched: false,
			}

			err := repo.Create(storedAssignment)
			So(err, ShouldBeNil)

			Convey("When fetching an assignment with the same ID", func() {
				fetchedAssignment, err := repo.Fetch(storedAssignment.ID)
				So(err, ShouldBeNil)

				Convey("Then the fetched assignment should be equal to the created assignment", func() {
					So(fetchedAssignment, ShouldResemble, storedAssignment)
				})
			})

			Convey("When creating an assignment with the same ID", func() {
				storedAssignment.Dispatched = true
				err := repo.Create(storedAssignment)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("When updating the assignment", func() {
				updatedAssignment := storedAssignment
				updatedAssignment.Dispatched = true
				err := repo.Update(updatedAssignment)
				So(err, ShouldBeNil)

				Convey("Then fetching the assignment should return the updated assignment", func() {
					fetchedAssignment, err := repo.Fetch(storedAssignment.ID)
					So(err, ShouldBeNil)
					So(fetchedAssignment, ShouldResemble, updatedAssignment)
				})
			})

			Convey("When deleting the assignment", func() {
				err := repo.Delete(storedAssignment.ID)
				So(err, ShouldBeNil)

				Convey("Then fetching the assignment should give a zero assignment and an error", func() {
					fetchedAssignment, err := repo.Fetch(storedAssignment.ID)
					So(fetchedAssignment, ShouldBeZeroValue)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

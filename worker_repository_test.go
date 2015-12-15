package junction

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"
)

func TestPostgresWorkerRepository(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer db.Close()
	defer db.Exec("DELETE FROM junction.workers")
	pwr := &PostgresWorkerRepository{db: db}

	runWorkerRepositoryTests(t, pwr)
}

func TestMapWorkerRepository(t *testing.T) {
	repo := &MapWorkerRepository{workers: make(map[string]Worker)}

	runWorkerRepositoryTests(t, repo)
}

func runWorkerRepositoryTests(t *testing.T, repo WorkerRepository) {
	Convey("Given a worker that has been created", t, func() {
		storedWorker := Worker{
			ID:            uuid.NewRandom(),
			Queue:         "test-queue",
			LastHeartbeat: nil,
			MaxJobCount:   10,
		}

		err := repo.Create(storedWorker)
		So(err, ShouldBeNil)

		Convey("When fetching a worker with the same ID", func() {
			fetchedWorker, err := repo.Fetch(storedWorker.ID)
			So(err, ShouldBeNil)

			Convey("Then the fetched worker should be equal to the stored worker", func() {
				So(fetchedWorker, ShouldResemble, storedWorker)
			})
		})

		Convey("When creating a worker with the same ID", func() {
			storedWorker.MaxJobCount = 15
			err := repo.Create(storedWorker)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When updating the worker", func() {
			storedWorker.MaxJobCount = 15
			err := repo.Update(storedWorker)
			So(err, ShouldBeNil)

			Convey("Then fetching the worker should return the updated attributes", func() {
				fetchedWorker, err := repo.Fetch(storedWorker.ID)
				So(err, ShouldBeNil)
				So(fetchedWorker, ShouldResemble, storedWorker)
			})
		})

		Convey("When deleting the worker", func() {
			err := repo.Delete(storedWorker.ID)
			So(err, ShouldBeNil)

			Convey("Then attempting to fetch the worker again should give a zero worker and an error", func() {
				fetchedWorker, err := repo.Fetch(storedWorker.ID)
				So(fetchedWorker, ShouldBeZeroValue)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a worker with a last-heartbeat timestamp that has been stored", t, func() {
		now := time.Now()
		storedWorker := Worker{
			ID:            uuid.NewRandom(),
			Queue:         "test-queue",
			LastHeartbeat: &now,
			MaxJobCount:   10,
		}

		err := repo.Create(storedWorker)
		So(err, ShouldBeNil)

		Convey("When fetching the worker", func() {
			fetchedWorker, err := repo.Fetch(storedWorker.ID)
			So(err, ShouldBeNil)

			Convey("Then the timestamp should be within a second of the stored timestamp", func() {
				So(fetchedWorker.LastHeartbeat, ShouldNotBeNil)
				So(*fetchedWorker.LastHeartbeat, ShouldHappenWithin, time.Second, *storedWorker.LastHeartbeat)
			})
		})
	})
}

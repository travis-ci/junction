package junction

import (
	"database/sql"
	"log"
	"os"
)

func Main(cfg *Config) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	workerRepo := &PostgresWorkerRepository{db: db}

	srv := newServer(cfg.Addr, workerRepo)

	srv.Setup()
	srv.Run()
}

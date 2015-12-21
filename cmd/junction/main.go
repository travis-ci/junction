package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/travis-ci/junction/database"
	junctionhttp "github.com/travis-ci/junction/http"
	"github.com/travis-ci/junction/junction"
)

func main() {
	app := cli.NewApp()
	app.Name = "junction"
	app.Usage = "Start the Junction HTTP server"
	app.Action = runJunction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Usage: "TCP address to listen on",
			Value: func() string {
				v := ":" + os.Getenv("PORT")
				if v == ":" {
					// Bind to a random port
					v = ":0"
				}
				return v
			}(),
			EnvVar: "JUNCTION_ADDR",
		},
		cli.StringSliceFlag{
			Name:   "worker-token",
			Usage:  "List of tokens to use for workers",
			EnvVar: "JUNCTION_WORKER_TOKENS",
		},
		cli.StringFlag{
			Name:   "database-url",
			Usage:  "URL to Postgres database to connect to",
			EnvVar: "JUNCTION_DATABASE_URL,DATABASE_URL",
		},
		cli.IntFlag{
			Name:   "database-max-pool-size",
			Usage:  "The maximum number of open connection to keep for the Postgres database",
			Value:  10,
			EnvVar: "JUNCTION_DATABASE_MAX_POOL_SIZE",
		},
	}

	app.Run(os.Args)
}

func runJunction(c *cli.Context) {
	database, err := database.NewPostgres(&database.PostgresConfig{
		URL:          c.String("database-url"),
		MaxOpenConns: c.Int("database-max-pool-size"),
	})
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	coreConfig := &junction.CoreConfig{
		Database:     database,
		WorkerTokens: c.StringSlice("worker-token"),
	}

	core, err := junction.NewCore(coreConfig)
	if err != nil {
		log.Fatalf("Error initializing core: %s", err)
	}

	server := &http.Server{
		Handler: junctionhttp.Handler(core),
	}

	listener, err := net.Listen("tcp", c.String("addr"))
	if err != nil {
		log.Fatalf("Error listening on TCP: %s", err)
	}

	log.Printf("Listening on %s", listener.Addr().String())

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Error serving on HTTP: %s", err)
	}
}

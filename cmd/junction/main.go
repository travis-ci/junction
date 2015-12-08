package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/travis-ci/junction"
)

func main() {
	app := cli.NewApp()
	app.Name = "junction"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "addr",
			Value: func() string {
				return ":" + os.Getenv("PORT")
			}(),
			Usage:  "host:port for HTTP server to bind to",
			EnvVar: "JUNCTION_ADDR,ADDR",
		},
		cli.StringFlag{
			Name:   "database-url",
			Usage:  "URL to Postgres database to connect to",
			EnvVar: "JUNCTION_DATABASE_URL,DATABASE_URL",
		},
	}

	app.Run(os.Args)
}

func action(c *cli.Context) {
	cfg := &junction.Config{
		Addr:        c.String("addr"),
		DatabaseURL: c.String("database-url"),
	}

	junction.Main(cfg)
}

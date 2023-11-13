package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	localhost = "http://localhost:3000/"
)

var (
	host string
)

func main() {
	startCmd := &cli.Command{
		Name:  "start",
		Usage: "Starts a cloudlife application",
		Action: func(cCtx *cli.Context) error {
			startMultiverse()
			return nil
		},
	}

	stopCmd := &cli.Command{
		Name:  "stop",
		Usage: "Stops a cloudlife application",
		Action: func(cCtx *cli.Context) error {
			stopMultiverse()
			return nil
		},
	}

	runCmd := &cli.Command{
		Name:  "run",
		Usage: "Runs a cloudlife application",
		Action: func(cCtx *cli.Context) error {
			runMultiverse()
			return nil
		},
	}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       localhost,
				Usage:       "Host to use to connect to the cloudlife application",
				Destination: &host,
			},
		},
		UsageText: "lifectl [global options] command [command options] [arguments]",
		Commands:  []*cli.Command{startCmd, runCmd, stopCmd},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

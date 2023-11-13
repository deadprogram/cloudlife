package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		Name:      "start",
		Usage:     "Starts a cloudlife application",
		UsageText: "lifectl start [size of multiverse]",
		Action: func(cCtx *cli.Context) error {
			size := 4
			if cCtx.NArg() > 0 {
				size, _ = strconv.Atoi(cCtx.Args().Get(0))
			}

			result, err := startMultiverse(size)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(result)
			return nil
		},
	}

	stopCmd := &cli.Command{
		Name:  "stop",
		Usage: "Stops a cloudlife application",
		Action: func(cCtx *cli.Context) error {
			result, err := stopMultiverse()
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println(result)
			return nil
		},
	}

	runCmd := &cli.Command{
		Name:  "run",
		Usage: "Runs the cloudlife application",
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
		Name:      "lifectl",
		Usage:     "CLI for cloudlife",
		UsageText: "lifectl [global options] command [command options] [arguments]",
		Commands:  []*cli.Command{startCmd, runCmd, stopCmd},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

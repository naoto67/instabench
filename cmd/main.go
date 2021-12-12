package main

import (
	"log"
	"os"

	"github.com/naoto67/instabench/internal"
	"github.com/urfave/cli/v2"
)

const version = "1.0.0"

func main() {
	app := &cli.App{
		Name:    "instabench",
		Usage:   "",
		Version: version,
		Commands: []*cli.Command{
			internal.NewBenchCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

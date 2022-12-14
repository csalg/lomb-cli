package main

import (
	"fmt"
	"os"

	"github.com/csalg/lomb-cli/pkg/commands/revise"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lomb",
		Usage: "Lomb CLI",
		Commands: []*cli.Command{
			revise.Cmd(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, "Error: "+err.Error())
	}
}

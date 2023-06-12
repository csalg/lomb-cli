package main

import (
	"fmt"
	"os"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/gpt"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process"
	"github.com/csalg/lomb-cli/cmd/lomb/config"
	"github.com/csalg/lomb-cli/pkg/revise"
	"github.com/urfave/cli/v2"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}
	// TODO start a goroutine with a channel listening for revision / assisted reading events
	app := &cli.App{
		Name:  "lomb",
		Usage: "Lomb CLI",
		Commands: []*cli.Command{
			revise.Cmd(),
			gpt.Cmd(conf),
			process.Cmd(conf),
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, "Error: "+err.Error())
	}
}

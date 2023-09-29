package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/drill"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/gpt"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/process"
	"github.com/csalg/lomb-cli/pkg/revise"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	deps, conf := bootstrap.Bootstrap(&templateFS)
	// TODO start a goroutine with a channel listening for revision / assisted reading events
	app := &cli.App{
		Name:  "lomb",
		Usage: "Lomb CLI",
		Commands: []*cli.Command{
			revise.Cmd(),
			gpt.Cmd(conf),
			process.Cmd(conf),
			drill.Cmd(deps, conf),
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, "Error: "+err.Error())
	}
}

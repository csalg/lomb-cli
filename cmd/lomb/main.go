package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/gpt"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/study"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	deps, conf := bootstrap.Bootstrap(&templateFS)
	app := &cli.App{
		Name:  "lomb",
		Usage: "Lomb CLI",
		Commands: []*cli.Command{
			gpt.Cmd(conf),
			text.Cmd(conf),
			study.Cmd(deps, conf),
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, "Error: "+err.Error())
	}
}

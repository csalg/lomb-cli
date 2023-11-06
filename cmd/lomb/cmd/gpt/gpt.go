package gpt

import (
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/urfave/cli/v2"
)

func Cmd(conf types.Config) *cli.Command {
	return &cli.Command{
		Name:  "gpt",
		Usage: "gpt related commands",
		Subcommands: []*cli.Command{
			GenerateCmd(conf),
			SimplifyCmd(conf),
		},
	}
}

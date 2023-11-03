package text

import (
	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/urfave/cli/v2"
)

func Cmd(conf bootstrap.Config) *cli.Command {
	return &cli.Command{
		Name:  "text",
		Usage: "process text",
		Subcommands: []*cli.Command{
			ProcessCmd(conf),
			ReverseCmd(conf),
		},
	}
}

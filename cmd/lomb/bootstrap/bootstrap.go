package bootstrap

import (
	"embed"

	"github.com/csalg/lomb-cli/pkg/repository"
	"github.com/csalg/lomb-cli/pkg/types"
)

func Bootstrap(fs *embed.FS) (Dependencies, types.Config) {
	conf, err := repository.ReadConfig()
	if err != nil {
		panic("reading config: " + err.Error())
	}
	deps := Dependencies{
		FS: fs,
	}
	return deps, conf
}

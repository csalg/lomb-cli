package bootstrap

import "embed"

func Bootstrap(fs *embed.FS) (Dependencies, Config) {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic("getting current working directory: " + err.Error())
	// }
	conf, err := ReadConfig()
	if err != nil {
		panic("reading config: " + err.Error())
	}
	deps := Dependencies{
		FS: fs,
	}
	return deps, conf
}

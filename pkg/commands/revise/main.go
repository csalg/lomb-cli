package revise

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/csalg/lomb-cli/pkg/io"
	"github.com/csalg/lomb-cli/pkg/service"
	"github.com/urfave/cli/v2"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:  "revise",
		Usage: "Revise a list of unranked words",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "source-language", Aliases: []string{"s"}, Required: true},
			&cli.StringFlag{Name: "target-language", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "lemmas-file", Aliases: []string{"l"}, Required: true},
		},
		Action: func(ctx *cli.Context) error {
			sourceLang := ctx.String("source-language")
			targetLang := ctx.String("target-language")
			lemmasFilename := ctx.String("lemmas-file")
			// TODO: Serve template.
			config, err := io.ReadConfig()
			if err != nil {
				log.Fatal(err)
			}
			s := service.New(config)

			lemmas := io.ReadLemmas(lemmasFilename)
			for _, lemma := range lemmas {
				url, err := s.GetDictionaryURL(&service.GetDictionaryURLInput{
					SourceLanguage: sourceLang,
					TargetLanguage: targetLang,
					Lemma:          lemma,
				})
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(url)
			}
			return nil
		},
	}
}

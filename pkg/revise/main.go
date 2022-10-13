package revise

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

			lemmas := readLemmas(lemmasFilename)
			urlMap := newDictionaryURLMap(sourceLang, targetLang)
			for _, lemma := range lemmas {
				fmt.Println(urlMap(lemma))
			}
			return nil
		},
	}
}

func readLemmas(filename string) (lemmas []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	  lemma := scanner.Text()
	  if lemma == "" {
	    continue
	  }
		lemmas = append(lemmas, lemma)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lemmas
}

type dictionaryUrlMap func(string) string

func newDictionaryURLMap(sourceLang, targetLang string) dictionaryUrlMap {
	return func(lemma string) string {
		switch sourceLang {
		case "en":
			switch targetLang {
			case "en":
				return fmt.Sprintf("https://www.merriam-webster.com/dictionary/%s", lemma)
			}
		case "da":
			switch targetLang {
			case "da":
				return fmt.Sprintf("https://ordnet.dk/ddo/ordbog?query=%s", lemma)
			case "en":
				return fmt.Sprintf("https://en.bab.la/dictionary/danish-english/%s?trending=1", lemma)

			}
		case "es":
			switch targetLang {
			case "es":
				return fmt.Sprintf("https://dle.rae.es/%s?m=form", lemma)
			}
		}
		return ""
	}
}

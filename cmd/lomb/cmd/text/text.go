package text

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/urfave/cli/v2"
)

func Cmd(conf bootstrap.Config) *cli.Command {
	return &cli.Command{
		Name:  "process",
		Usage: "lemmatize and translate a text so it can be used for assisted reading or revision",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "target-language", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "base-language", Aliases: []string{"b"}, Required: true},
			&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Required: true},
		},
		Action: func(ctx *cli.Context) error {
			text, err := readText(ctx.String("file"))
			if err != nil {
				return fmt.Errorf("reading text: %w", err)
			}
			baseLang, err := types.NewLanguage(ctx.String("base-language"))
			if err != nil {
				return fmt.Errorf("parsing base language: %w", err)
			}
			targetLang, err := types.NewLanguage(ctx.String("target-language"))
			if err != nil {
				return fmt.Errorf("parsing target language: %w", err)
			}
			tp, err := NewTextProcessor(Config{
				BaseLanguage:   baseLang,
				TargetLanguage: targetLang,
				DeeplAPIKey:    conf.DeeplAPIKey,
				DeeplAPIPro:    conf.DeeplAPIPro,
			})
			if err != nil {
				return fmt.Errorf("creating text processor: %w", err)
			}
			processedText, err := tp.Process(text)
			if err != nil {
				return fmt.Errorf("processing text: %w", err)
			}
			if err := writeProcessedText(ctx.String("file")+".lotxt", processedText); err != nil {
				return fmt.Errorf("writing lotxt: %w", err)
			}
			if err := os.WriteFile(ctx.String("file")+".translated.txt", []byte(processedText.Translation()), 0o644); err != nil {
				return fmt.Errorf("writing translated text: %w", err)
			}
			return nil
		},
	}
}

func readText(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	return string(b), nil
}

func writeProcessedText(filename string, processedText types.ProcessedText) error {
	file, _ := json.MarshalIndent(processedText, "", " ")
	//nolint: gosec
	if err := os.WriteFile(filename, file, 0o644); err != nil {
		return err
	}
	return nil
}

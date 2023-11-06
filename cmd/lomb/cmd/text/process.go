package text

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/urfave/cli/v2"
)

func ProcessCmd(conf types.Config) *cli.Command {
	return &cli.Command{
		Name:  "process",
		Usage: "lemmatize and translate a text so it can be used for assisted reading or revision",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			text, err := readString(ctx.Args().First())
			if err != nil {
				return fmt.Errorf("reading text: %w", err)
			}
			tp, err := NewTextProcessor(Config{
				BaseLanguage:   conf.BaseLanguage,
				TargetLanguage: conf.TargetLanguage,
				DeeplAPIKey:    conf.DeeplAPIKey,
				DeeplAPIPro:    conf.DeeplAPIPro,
				OpenAIAPIKey:   conf.OpenAIAPIKey,
			})
			if err != nil {
				return fmt.Errorf("creating text processor: %w", err)
			}
			processedText, err := tp.Process(text, ProcessOptions{})
			if err != nil {
				return fmt.Errorf("processing text: %w", err)
			}
			if err := writeJSON(ctx.String("file")+".lotxt", processedText); err != nil {
				return fmt.Errorf("writing lotxt: %w", err)
			}
			if err := os.WriteFile(ctx.String("file")+".translated.txt", []byte(processedText.Translation()), 0o644); err != nil {
				return fmt.Errorf("writing translated text: %w", err)
			}
			return nil
		},
	}
}

func readString(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	return string(b), nil
}

func writeJSON(filename string, i any) error {
	file, _ := json.MarshalIndent(i, "", " ")
	//nolint: gosec
	if err := os.WriteFile(filename, file, 0o644); err != nil {
		return err
	}
	return nil
}

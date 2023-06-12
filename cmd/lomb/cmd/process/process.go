package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/config"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/urfave/cli/v2"
)

func Cmd(conf config.Config) *cli.Command {
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
			processedText, err := processText(text, ctx.String("base-language"), ctx.String("target-language"))
			if err != nil {
				return fmt.Errorf("processing text: %w", err)
			}
			if err := writeProcessedText(ctx.String("file")+".lotxt", processedText); err != nil {
				return fmt.Errorf("writing processed text: %w", err)
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

func processText(text, baseLanguage, targetLanguage string) (types.ProcessedText, error) {
	processedText := types.ProcessedText{
		BaseLanguage:   baseLanguage,
		TargetLanguage: targetLanguage,
	}
	paragraphsStrArr := strings.Split(text, "\n")
	for _, paragraphStr := range paragraphsStrArr {
		paragraph := types.Paragraph{}
		sentencesStrArr := strings.Split(paragraphStr, ".")
		for _, sentenceStr := range sentencesStrArr {
			chunk := types.Chunk{
				Translation: sentenceStr,
			}
			for _, word := range strings.Split(sentenceStr, " ") {
				chunk.Tokens = append(chunk.Tokens, types.Token{
					Text:  word,
					Lemma: word,
				})
			}
			paragraph = append(paragraph, chunk)
		}
		processedText.Paragraphs = append(processedText.Paragraphs, paragraph)
	}
	return processedText, nil
}

func writeProcessedText(filename string, processedText types.ProcessedText) error {
	file, _ := json.MarshalIndent(processedText, "", " ")
	//nolint: gosec
	if err := os.WriteFile(filename, file, 0o644); err != nil {
		return err
	}
	return nil
}

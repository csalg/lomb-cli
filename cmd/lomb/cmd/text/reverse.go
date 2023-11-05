package text

import (
	"fmt"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/lemmatizers/dummylemmatizer"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func ReverseCmd(conf types.Config) *cli.Command {
	return &cli.Command{
		Name:  "reverse",
		Usage: "reverse the languages of a processed text",
		Action: func(ctx *cli.Context) error {
			filename := ctx.Args().First()
			if filename == "" {
				return cli.Exit("filename is required", 1)
			}
			var txt types.ProcessedText
			found, err := utils.ReadAndUnmarshal(filename, &txt)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			if !found {
				return cli.Exit("file not found", 1)
			}

			lemmatizer := dummylemmatizer.New()
			reversedTxt := types.ProcessedText{}

			for _, p := range txt.Paragraphs {
				reversedParagraph := types.Paragraph{}
				for _, s := range p {
					tokens, err := lemmatizer.Lemmatize(s.Translation)
					if err != nil {
						return fmt.Errorf("lemmatizing: %w", err)
					}
					reversedParagraph = append(reversedParagraph, types.Chunk{
						Tokens:      tokens,
						Translation: s.Sentence(),
					})
				}
				reversedTxt.Paragraphs = append(reversedTxt.Paragraphs, reversedParagraph)
			}
			if err := writeProcessedText(filename+"-reversed", reversedTxt); err != nil {
				return fmt.Errorf("writing lotxt: %w", err)
			}
			return nil
		},
	}
}

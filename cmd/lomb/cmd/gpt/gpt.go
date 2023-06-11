package gpt

import (
	"fmt"
	"os"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/gpt/openai"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// The way this will work is that the user will pass a filename with yaml.
// The yaml has to be converted to a series of prompts.
//
// Then each prompt is passed to gpt with the previous user message and the gpt response.
// The gpt response is then parsed and the next prompt is generated and so on.

// This a sample of the yaml to parse

// For now, just call the api and ask it to reply hello world.

func Cmd() *cli.Command {
	return &cli.Command{
		Name:  "gpt",
		Usage: "generate text with gpt",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Required: true},
		},
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			prompts, err := parseYaml(filename)
			if err != nil {
				return fmt.Errorf("parsing yaml: %w", err)
			}

			prompt := prompts[0]
			res, err := openai.Post([]openai.Message{
				{Content: prompt, Role: openai.MessageRoleUser},
			}, "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")
			if err != nil {
				return fmt.Errorf("calling openai: %w", err)
			}
			fmt.Println(res)
			return nil
		},
	}
}

type PromptDoc struct {
	Language string   `yaml:"language"`
	Prompt   string   `yaml:"prompt"`
	Topics   []string `yaml:"topics"`
}

func parseYaml(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var doc PromptDoc
	if err := yaml.NewDecoder(f).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decoding yaml: %w", err)
	}
	basePrompt := strings.ReplaceAll(doc.Prompt, "$TARGET_LANGUAGE", doc.Language)
	if len(doc.Topics) == 0 {
		return []string{basePrompt}, nil
	}
	var prompts []string
	for _, topic := range doc.Topics {
		prompt := strings.ReplaceAll(basePrompt, "$TOPIC", topic)
		prompts = append(prompts, prompt)
	}
	return prompts, nil
}

package gpt

import (
	"fmt"
	"os"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/cmd/lomb/cmd/gpt/openai"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const costPer1000Tokens = 0.002

func Cmd(conf bootstrap.Config) *cli.Command {
	return &cli.Command{
		Name:  "gpt",
		Usage: "generate text with gpt",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Required: true},
		},
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			prompts, err := extractPrompts(filename)
			if err != nil {
				return fmt.Errorf("parsing yaml: %w", err)
			}

			responses := []string{}
			totalCost := 0.0
			for i, prompt := range prompts {
				fmt.Printf("Prompt %d/%d \n", i+1, len(prompts))
				res, err := openai.Post([]openai.Message{
					{Content: prompt, Role: openai.MessageRoleUser},
				}, conf.OpenAIAPIKey)
				if err != nil {
					return fmt.Errorf("calling openai: %w", err)
				}
				tokens := res.Usage.TotalTokens
				cost := float64(tokens) / 1000 * costPer1000Tokens
				fmt.Printf("Tokens: %d, cost: $%.4f\n", tokens, cost)
				totalCost += cost
				fmt.Printf("Total cost: $%.4f\n", totalCost)
				firstChoice := res.Choices[0].Message.Content
				fmt.Printf("Response:\n%s\n", firstChoice)
				responses = append(responses, firstChoice)
				fmt.Println()
			}
			if err := writeResponses(filename+".txt", responses); err != nil {
				return fmt.Errorf("writing responses: %w", err)
			}
			return nil
		},
	}
}

type PromptDoc struct {
	Language string  `yaml:"language"`
	Prompt   string  `yaml:"prompt"`
	Topics   []Topic `yaml:"topics"`
}

type Topic struct {
	Topic string `yaml:"topic"`
}

func extractPrompts(filename string) ([]string, error) {
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
		prompt := strings.ReplaceAll(basePrompt, "$TOPIC", topic.Topic)
		prompts = append(prompts, prompt)
	}
	return prompts, nil
}

func writeResponses(filename string, responses []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	resString := strings.Join(responses, "\n----\n")
	if _, err := f.WriteString(resString); err != nil {
		return fmt.Errorf("writing responses: %w", err)
	}
	return nil
}

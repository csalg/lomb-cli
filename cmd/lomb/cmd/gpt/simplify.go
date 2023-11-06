package gpt

import (
	"fmt"
	"os"
	"strings"

	"github.com/csalg/lomb-cli/pkg/openai"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/urfave/cli/v2"
)

func SimplifyCmd(conf types.Config) *cli.Command {
	return &cli.Command{
		Name:  "simplify",
		Usage: "simplify text with gpt",
		Action: func(ctx *cli.Context) error {
			text, err := os.ReadFile(ctx.Args().First())
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}

			newText := ""
			totalCost := 0.0
			for _, line := range strings.Split(string(text), "\n") {
				if line == "" {
					continue
				}

				prompt := fmt.Sprintf(`Rewrite this text so it is suitable for someone who is learning English. SIMPLIFY the text.

				Use simple words and short sentences. It is better to say "He is sad" than "He was terribly depressed".
				Use simple grammar.
				You can write many simple sentences instead of one complex sentence.
				Make sure every detail is clear.
				Make sure the name of everything is clear.
				Repeat many times if necessary.
				Explain what you mean if necessary.

				It must be easy to understand the text, even if you just started learning English.

				Here is the text:
				%s
				`, line)
				if len(line) < 20 {
					newText += "\n" + line
					continue
				}

				res, err := openai.Post([]openai.Message{
					{Content: prompt, Role: openai.MessageRoleUser},
				}, conf.OpenAIAPIKey)
				if err != nil {
					return fmt.Errorf("calling openai: %w", err)
				}
				tokens := res.Usage.TotalTokens
				cost := float64(tokens) / 1000 * openai.CostPer1000Tokens
				fmt.Printf("Tokens: %d, cost: $%.4f\n", tokens, cost)
				totalCost += cost
				fmt.Printf("Total cost: $%.4f\n", totalCost)
				firstChoice := res.Choices[0].Message.Content
				fmt.Printf("Original text:\n%s\n", line)
				fmt.Printf("Response:\n%s\n", firstChoice)
				newText += "\n" + firstChoice
				fmt.Println()
			}

			basename := strings.TrimSuffix(ctx.Args().First(), ".txt")
			newFilename := basename + "-simplified.txt"
			if err := os.WriteFile(newFilename, []byte(newText), 0644); err != nil {
				return fmt.Errorf("writing file: %w", err)
			}
			return nil
		},
	}
}

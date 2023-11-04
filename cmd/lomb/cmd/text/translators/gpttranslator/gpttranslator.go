package gpttranslator

import (
	"fmt"
	"strings"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/text/translators"
	"github.com/csalg/lomb-cli/pkg/openai"
	"github.com/csalg/lomb-cli/pkg/types"
)

// type Translator interface {
// 	Translate(sourceLang, targetLang types.Language, text []string) ([]translators.TranslatedText, error)
// }

type GPT struct {
	OpenAIAPIKey string
}

func New(openAIAPIKey string) *GPT {
	return &GPT{
		OpenAIAPIKey: openAIAPIKey,
	}
}

// Translate translates a list of texts from a source language to a target language.
func (g *GPT) Translate(sourceLang, targetLang types.Language, texts []string) ([]translators.TranslatedText, error) {
	translatedTexts := []translators.TranslatedText{}
	totalCost := 0.0
	for _, text := range texts {
		if len(sanitize(text)) == 0 {
			translatedTexts = append(translatedTexts, translators.TranslatedText{
				Original:   text,
				Translated: text,
			})
			continue
		}
		fmt.Printf("Translating text: %q\n", text)
		prompt := fmt.Sprintf(`Translate from '%s' to '%s'. Output ONLY the translation.`,
			sourceLang, targetLang)
		prompt += "\n"
		prompt += text
		fmt.Printf("Prompt:\n%s\n", prompt)

		res, err := openai.Post([]openai.Message{
			{Content: prompt, Role: openai.MessageRoleUser},
		}, g.OpenAIAPIKey)
		if err != nil {
			return []translators.TranslatedText{}, fmt.Errorf("calling openai: %w", err)
		}
		tokens := res.Usage.TotalTokens
		cost := float64(tokens) / 1000 * openai.CostPer1000Tokens
		fmt.Printf("Tokens: %d, cost: $%.4f\n", tokens, cost)
		totalCost += cost
		fmt.Printf("Total cost: $%.4f\n", totalCost)
		firstChoice := res.Choices[0].Message.Content
		fmt.Printf("Response:\n%v\n", firstChoice)
		translatedTexts = append(translatedTexts, translators.TranslatedText{
			Original:   text,
			Translated: firstChoice,
		})
		fmt.Println()
	}
	return translatedTexts, nil
}

// sanitize
func sanitize(text string) string {
	sanitized := strings.TrimSpace(text)
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")
	return sanitized
}

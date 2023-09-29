package drill

import (
	"fmt"
	"sort"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func Cmd(deps bootstrap.Dependencies, conf bootstrap.Config) *cli.Command {
	return &cli.Command{
		Name:  "drill",
		Usage: "drill vocabulary from lotxt files",
		Action: func(ctx *cli.Context) error {
			filename := ctx.Args().First()
			if filename == "" {
				return cli.Exit("filename is required", 1)
			}
			fmt.Println(filename)
			var txt types.ProcessedText
			found, err := utils.ReadAndUnmarshal(filename, &txt)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			if !found {
				return cli.Exit("file not found", 1)
			}
			lemmaCounts := getLemmaCounts(txt)
			fmt.Println(lemmaCounts)
			return nil
		},
	}
}

func getLemmaCounts(txt types.ProcessedText) []LemmaCount {
	lemmaCountLookup := make(map[string]int)
	for _, paragraph := range txt.Paragraphs {
		for _, chunk := range paragraph {
			for _, token := range chunk.Tokens {
				lemmaCountLookup[token.Lemma]++
			}
		}
	}
	var lemmaCounts []LemmaCount
	for lemma, count := range lemmaCountLookup {
		lemmaCounts = append(lemmaCounts, LemmaCount{Lemma: lemma, Count: count})
	}
	sort.Slice(lemmaCounts, func(i, j int) bool {
		return lemmaCounts[i].Count > lemmaCounts[j].Count
	})
	return lemmaCounts
}

type LemmaCount struct {
	Lemma string
	Count int
}

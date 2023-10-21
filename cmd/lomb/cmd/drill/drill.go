package drill

import (
	"embed"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	numbersRegex         = regexp.MustCompile(`\d+`)
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
			var txt types.ProcessedText
			found, err := utils.ReadAndUnmarshal(filename, &txt)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			if !found {
				return cli.Exit("file not found", 1)
			}
			driller := NewDriller(deps.FS, Config{
				Port:      conf.Port,
				ClientURL: conf.ClientURL(),
			})
			driller.LoadLemmaCounts(txt)
			driller.OpenURLInBrowser()
			driller.Serve()
			return nil
		},
	}
}

type Driller struct {
	FS          *embed.FS
	lemmaCounts []LemmaCount
	config      Config
}

type LemmaCount struct {
	Lemma string
	Count int
}

type Config struct {
	Port      string
	ClientURL string
}

type PageModel struct {
	LemmaCounts []LemmaCount
}

func NewDriller(fs *embed.FS, conf Config) *Driller {
	return &Driller{
		FS:     fs,
		config: conf,
	}
}

// LoadLemmaCounts loads the lemma counts from the given texts.
func (d *Driller) LoadLemmaCounts(txts ...types.ProcessedText) {
	lemmaCountLookup := make(map[string]int)
	for _, txt := range txts {
		for _, paragraph := range txt.Paragraphs {
			for _, chunk := range paragraph {
				for _, token := range chunk.Tokens {
					// Trim alphanumeric characters from the lemma.
					sanitizedLemma := sanitizeLemma(token.Lemma)
					if sanitizedLemma == "" {
						continue
					}
					lemmaCountLookup[sanitizedLemma]++
				}
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
	d.lemmaCounts = lemmaCounts
}

// OpenURLInBrowser opens the given URL in the default browser of the user.
func (d *Driller) OpenURLInBrowser() {
	go utils.OpenURLInBrowser(d.config.ClientURL)
}

// Serve starts the web server.
func (d *Driller) Serve() {
	mux := chi.NewRouter()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, r, d.FS, "templates/drill.html", PageModel{
			LemmaCounts: d.lemmaCounts,
		})
	})
	s := &http.Server{
		Addr:              ":" + d.config.Port,
		Handler:           mux,
		ReadTimeout:       100 * time.Second,
		ReadHeaderTimeout: 100 * time.Second,
		WriteTimeout:      100 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	panic(s.ListenAndServe())
}

func sanitizeLemma(str string) string {
	cleanStr := nonAlphanumericRegex.ReplaceAllString(str, "")
	cleanStr = numbersRegex.ReplaceAllString(cleanStr, "")
	return strings.TrimSpace(strings.ToLower(cleanStr))
}

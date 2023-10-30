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
	jsoniter "github.com/json-iterator/go"
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
	FS                 *embed.FS
	config             Config
	lemmaCounts        []LemmaCount
	lemmaExampleLookup map[string][]Example
}

type Config struct {
	Port      string
	ClientURL string
}

func NewDriller(fs *embed.FS, conf Config) *Driller {
	return &Driller{
		FS:                 fs,
		config:             conf,
		lemmaExampleLookup: make(map[string][]Example),
	}
}

// LoadLemmaCounts loads the lemma counts from the given texts.
func (d *Driller) LoadLemmaCounts(txts ...types.ProcessedText) {
	lemmaExampleLookup := make(map[string][]Example)
	for _, txt := range txts {
		for _, paragraph := range txt.Paragraphs {
			for _, chunk := range paragraph {
				sentence := Example{
					Text:        chunk.Sentence(),
					Translation: chunk.Translation,
				}
				for _, token := range chunk.Tokens {
					// Trim alphanumeric characters from the lemma.
					sanitizedLemma := sanitizeLemma(token.Lemma)
					if sanitizedLemma == "" {
						continue
					}
					lemmaExampleLookup[sanitizedLemma] = append(lemmaExampleLookup[sanitizedLemma], sentence)
				}
			}
		}
	}
	var lemmaCounts []LemmaCount
	for lemma, examples := range lemmaExampleLookup {
		lemmaCounts = append(lemmaCounts, LemmaCount{Lemma: lemma, Count: len(examples)})
	}
	sort.Slice(lemmaCounts, func(i, j int) bool {
		return lemmaCounts[i].Count > lemmaCounts[j].Count
	})
	d.lemmaCounts = lemmaCounts
	d.lemmaExampleLookup = lemmaExampleLookup
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
			Grid: DefaultGrid,
			Data: Data{
				LemmaCounts: d.lemmaCounts,
			},
		})
	})

	mux.Post("/examples", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Lemma string `json:"lemma"`
		}
		decoder := utils.CreateJSONDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := jsoniter.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusBadRequest, Message: "invalid request"})
			if err != nil {
				fmt.Println("couldn't write body:", err)
			}
			return
		}
		examples, ok := d.lemmaExampleLookup[req.Lemma]
		examplesNoDuplicates := make([]Example, 0, len(examples))
		seen := make(map[string]struct{})
		for _, example := range examples {
			if _, ok := seen[example.Text]; ok {
				continue
			}
			seen[example.Text] = struct{}{}
			examplesNoDuplicates = append(examplesNoDuplicates, example)
		}
		// Sort by length
		sort.Slice(examplesNoDuplicates, func(i, j int) bool {
			return len(examplesNoDuplicates[i].Text) < len(examplesNoDuplicates[j].Text)
		})

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			err := jsoniter.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusNotFound, Message: "not found"})
			if err != nil {
				fmt.Println("couldn't write body:", err)
			}
			return
		}
		err := jsoniter.NewEncoder(w).Encode(struct {
			Translation []Example `json:"examples"`
		}{Translation: examplesNoDuplicates})
		if err != nil {
			fmt.Println("couldn't write body:", err)
		}
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

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

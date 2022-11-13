package revise

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/csalg/lomb-cli/pkg/io"
	"github.com/csalg/lomb-cli/pkg/lib"
	"github.com/csalg/lomb-cli/pkg/service"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
)

// RevisableLemma is a lemma that can be revised.
type RevisableLemma struct {
	Lemma          string
	URL            string
	SourceLanguage string
	TargetLanguage string
}

func Cmd() *cli.Command {
	return &cli.Command{
		Name:  "revise",
		Usage: "Revise a list of unranked words",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "source-language", Aliases: []string{"s"}, Required: true},
			&cli.StringFlag{Name: "target-language", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "lemmas-file", Aliases: []string{"l"}, Required: true},
			&cli.BoolFlag{Name: "stdout", Aliases: []string{"o"}, Required: false},
		},
		Action: func(ctx *cli.Context) error {
			sourceLang := ctx.String("source-language")
			targetLang := ctx.String("target-language")
			lemmasFilename := ctx.String("lemmas-file")
			// TODO: Serve template.
			config, err := io.ReadConfig()
			if err != nil {
				log.Fatal(err)
			}
			s := service.New(config)

			lemmas := io.ReadLemmas(lemmasFilename)
			revisableLemmas := []RevisableLemma{}
			for _, lemma := range lemmas {
				url, err := s.GetDictionaryURL(&service.GetDictionaryURLInput{
					SourceLanguage: sourceLang,
					TargetLanguage: targetLang,
					Lemma:          lemma,
				})
				if err != nil {
					log.Fatal(err)
				}
				revisableLemmas = append(revisableLemmas, RevisableLemma{
					Lemma:          lemma,
					URL:            url,
					SourceLanguage: sourceLang,
					TargetLanguage: targetLang,
				})
			}
			if ctx.Bool("stdout") {
				ToSTDOUT(revisableLemmas)
				return nil
			}
			ToServer(&config, revisableLemmas)
			return nil
		},
	}
}

// ToSTDOUT prints out the urls to STDOUT.
func ToSTDOUT(revisableLemmas []RevisableLemma) {
	for _, revisableLemma := range revisableLemmas {
		fmt.Println(revisableLemma.URL)
	}
}

// ToServer starts a server with the revision view.
func ToServer(conf *types.Config, revisableLemmas []RevisableLemma) {
	router := chi.NewRouter()
	router.Get("/", revisionHandler(conf, revisableLemmas))

	clientURL := fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)

	go lib.OpenURLInBrowser(clientURL)
	fmt.Println("Rofl Web UI is running on " + clientURL)
	s := &http.Server{
		Addr:              ":" + fmt.Sprint(conf.Port),
		Handler:           router,
		ReadTimeout:       100 * time.Second,
		ReadHeaderTimeout: 100 * time.Second,
		WriteTimeout:      100 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	panic(s.ListenAndServe())
}

type RevisionViewModels struct {
	Lemmas []RevisableLemma
}

// revisionHandler handles the revision view.
func revisionHandler(conf *types.Config, revisableLemmas []RevisableLemma) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := RevisionViewModels{
			Lemmas: revisableLemmas,
		}
		tmpl := template.Must(template.ParseFiles("templates/revision.html"))
		if err := tmpl.Execute(w, model); err != nil {
			panic("rendering template: " + err.Error())
		}
	}
}

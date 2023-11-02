package drill

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils"
	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
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
			var txt types.ProcessedText
			found, err := utils.ReadAndUnmarshal(filename, &txt)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			if !found {
				return cli.Exit("file not found", 1)
			}
			driller := NewServer(deps.FS, Config{
				Port:      conf.Port,
				ClientURL: conf.ClientURL(),
			})
			driller.corpus.LoadTexts(txt)
			driller.OpenURLInBrowser()
			driller.Serve()
			return nil
		},
	}
}

type Server struct {
	FS     *embed.FS
	config Config
	corpus *Corpus
}

type Config struct {
	Port      string
	ClientURL string
}

func NewServer(fs *embed.FS, conf Config) *Server {
	return &Server{
		FS:     fs,
		config: conf,
		corpus: NewCorpus(),
	}
}

// OpenURLInBrowser opens the given URL in the default browser of the user.
func (srv *Server) OpenURLInBrowser() {
	go utils.OpenURLInBrowser(srv.config.ClientURL)
}

// Serve starts the web server.
func (srv *Server) Serve() {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, r, srv.FS, "templates/drill.html", PageModel{
			Grid: AllDataViews,
			Data: Data{
				LemmaCounts:      srv.corpus.LemmaCounts,
				ReaderParagraphs: srv.corpus.Paragraphs,
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
		err := jsoniter.NewEncoder(w).Encode(struct {
			Examples []string `json:"examples"`
		}{Examples: srv.corpus.GetExamples(req.Lemma)})
		if err != nil {
			fmt.Println("couldn't write body:", err)
		}
	})

	mux.Post("/understandable-sentences", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			MinUnderstandability float64 `json:"min_understandability"`
			MaxLength            float64 `json:"max_length"`
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
		sentences := srv.corpus.GetUnderstandableSentences(req.MaxLength, req.MinUnderstandability)
		sentencesStr := make([]string, 0, len(sentences))
		for _, sentence := range sentences {
			sentencesStr = append(sentencesStr, sentence.Text)
		}

		err := jsoniter.NewEncoder(w).Encode(struct {
			Sentences []string `json:"sentences"`
		}{Sentences: sentencesStr})
		if err != nil {
			fmt.Println("couldn't write body:", err)
		}
	})

	mux.Post("/translate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Text string `json:"text"`
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
		err := jsoniter.NewEncoder(w).Encode(struct {
			Translation string `json:"translation"`
		}{Translation: srv.corpus.Translate(req.Text)})
		if err != nil {
			fmt.Println("couldn't write body:", err)
		}
	})

	s := &http.Server{
		Addr:              ":" + srv.config.Port,
		Handler:           mux,
		ReadTimeout:       100 * time.Second,
		ReadHeaderTimeout: 100 * time.Second,
		WriteTimeout:      100 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	panic(s.ListenAndServe())
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

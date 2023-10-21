package read

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
		Name:  "read",
		Usage: "read lotxt files",
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
			rs := NewReadServer(deps.FS, Config{
				Port:      conf.Port,
				ClientURL: conf.ClientURL(),
			})
			rs.OpenURLInBrowser()
			rs.Serve(txt)
			return nil
		},
	}
}

type ReadServer struct {
	FS     *embed.FS
	config Config
}

type Config struct {
	Port      string
	ClientURL string
}

// NewReadServer creates a new drill server.
func NewReadServer(fs *embed.FS, config Config) *ReadServer {
	return &ReadServer{
		FS:     fs,
		config: config,
	}
}

// OpenURLInBrowser opens the given URL in the default browser of the user.
func (rs *ReadServer) OpenURLInBrowser() {
	go utils.OpenURLInBrowser(rs.config.ClientURL)
}

// Serve starts the web server.
func (rs *ReadServer) Serve(txt types.ProcessedText) {
	mux := chi.NewRouter()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, r, rs.FS, "templates/read.html", txt)
	})
	translations := make(map[string]string)
	for _, p := range txt.Paragraphs {
		for _, c := range p {
			translations[c.Sentence()] = c.Translation
		}
	}
	mux.Post("/translate", func(w http.ResponseWriter, r *http.Request) {
		// Body has the form: {"sentence": "sentence to translate"}
		var req struct {
			Sentence string `json:"sentence"`
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
		translation, ok := translations[req.Sentence]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			err := jsoniter.NewEncoder(w).Encode(ErrorResponse{Code: http.StatusNotFound, Message: "not found"})
			if err != nil {
				fmt.Println("couldn't write body:", err)
			}
			return
		}
		err := jsoniter.NewEncoder(w).Encode(struct {
			Translation string `json:"translation"`
		}{Translation: translation})
		if err != nil {
			fmt.Println("couldn't write body:", err)
		}
	})
	s := &http.Server{
		Addr:              ":" + rs.config.Port,
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

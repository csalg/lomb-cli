package interact

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/csalg/lomb-cli/cmd/lomb/bootstrap"
	"github.com/csalg/lomb-cli/pkg/types"
	"github.com/csalg/lomb-cli/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"
)

func Cmd(deps bootstrap.Dependencies, conf types.Config) *cli.Command {
	return &cli.Command{
		Name:  "interact",
		Usage: "serve translation UI",
		Action: func(ctx *cli.Context) error {
			srv := NewServer(deps.FS, Config{
				Port:      conf.Port,
				ClientURL: conf.ClientURL(),
			})
			srv.OpenURLInBrowser()
			srv.Serve()
			return nil
		},
	}
}

type User struct {
	Name           string
	KnownLanguage  string
	TargetLanguage string
}

var users = []User{
	{
		Name:           "Carlos",
		KnownLanguage:  "English",
		TargetLanguage: "Bulgarian",
	},
	{
		Name:           "Ivan",
		KnownLanguage:  "Bulgarian",
		TargetLanguage: "English",
	},
}

type Server struct {
	FS     *embed.FS
	config Config
}

type Config struct {
	Port      string
	ClientURL string
}

func NewServer(fs *embed.FS, conf Config) *Server {
	return &Server{
		FS:     fs,
		config: conf,
	}
}

// OpenURLInBrowser opens the given URL in the default browser of the user.
func (srv *Server) OpenURLInBrowser() {
	go utils.OpenURLInBrowser(srv.config.ClientURL)
}

// Serve starts the web server.
func (srv *Server) Serve() {
	mux := chi.NewRouter()

	for _, user := range users {
		mux.Get("/"+strings.ToLower(user.Name), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", user.Name)
		})
	}
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

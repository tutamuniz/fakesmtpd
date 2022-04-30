package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/tutamuniz/fakesmtpd/internal/config"
)

//go:embed web/index.html
var content embed.FS

type Server struct {
	config *config.Config
}

func NewHTTPServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (h *Server) Start() {
	go h.processMessages()
}

func (h *Server) processMessages() {
	content, _ := fs.Sub(content, "web")

	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/capture/enable", func(wr http.ResponseWriter, r *http.Request) {
		h.config.EnableCapture()
		wr.Write([]byte("ENABLECAP OK"))
	})

	http.HandleFunc("/capture/disable", func(wr http.ResponseWriter, r *http.Request) {
		h.config.DisableCapture()
		wr.Write([]byte("DISABLECAP OK"))
	})

	http.HandleFunc("/capture/status", func(wr http.ResponseWriter, r *http.Request) {
		if h.config.CaptureStatus {
			wr.Write([]byte("CAPENABLED"))
		} else {
			wr.Write([]byte("CAPDISABLED"))
		}
	})

	http.ListenAndServe(":8080", nil)
}

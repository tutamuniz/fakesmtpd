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

	logger := h.config.Logger

	http.Handle("/", http.FileServer(http.FS(content)))

	http.HandleFunc("/capture/enable", h.enableHandler)

	http.HandleFunc("/capture/disable", h.disableHandler)

	http.HandleFunc("/capture/status", h.statusHandler)

	logger.Println("Starting HTTP server")
	logger.Fatal("%s", http.ListenAndServe(h.config.HTTPServerConfig.Address, nil))
}

func (h *Server) enableHandler(wr http.ResponseWriter, r *http.Request) {
	h.config.EnableCapture()
	_, err := wr.Write([]byte("ENABLECAP OK"))
	if err != nil {
		h.config.Logger.Printf("Error: %s", err)
	}
}

func (h *Server) disableHandler(wr http.ResponseWriter, r *http.Request) {
	h.config.DisableCapture()
	_, err := wr.Write([]byte("DISABLECAP OK"))
	if err != nil {
		h.config.Logger.Printf("Error: %s", err)
	}
}

func (h *Server) statusHandler(wr http.ResponseWriter, r *http.Request) {
	if h.config.CaptureStatus {
		_, err := wr.Write([]byte("CAPENABLED"))
		if err != nil {
			h.config.Logger.Printf("Error writing to response: %s", err)
		}

	} else {
		_, err := wr.Write([]byte("CAPDISABLED"))
		if err != nil {
			h.config.Logger.Printf("Error writing to response: %s", err)
		}

	}
}

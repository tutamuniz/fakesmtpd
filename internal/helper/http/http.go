package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/tutamuniz/fakesmtpd/internal/server"
)

//go:embed web/index.html
var content embed.FS

func Server(fake *server.FakeSMTP) { // terrible, but works
	content, _ := fs.Sub(content, "web")

	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/capture/enable", func(wr http.ResponseWriter, r *http.Request) {
		fake.EnableCapture()
		wr.Write([]byte("ENABLECAP OK"))
	})

	http.HandleFunc("/capture/disable", func(wr http.ResponseWriter, r *http.Request) {
		fake.DisableCapture()
		wr.Write([]byte("DISABLECAP OK"))
	})

	http.HandleFunc("/capture/status", func(wr http.ResponseWriter, r *http.Request) {
		status := fake.CaptureStatus()
		if status {
			wr.Write([]byte("CAPENABLED"))
		} else {
			wr.Write([]byte("CAPDISABLED"))
		}
	})

	http.ListenAndServe(":8080", nil)
}

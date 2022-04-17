package web

import (
	"io/fs"
	"net/http"

	"github.com/tutamuniz/fakesmtpd/server"
)

func Server(fake *server.FakeSMTP, content fs.FS) { // terrible, but works
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

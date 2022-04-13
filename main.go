package main

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/tutamuniz/fakesmtpd/helper/chat"
	"github.com/tutamuniz/fakesmtpd/server"
)

//go:embed embfs/index.html
var content embed.FS

func main() {
	root, _ := fs.Sub(content, "embfs")

	path := os.Args[1]

	fakeServer := server.NewServer("0.0.0.0:25", path)

	bot := chat.NewBot(fakeServer.Logger)

	fakeServer.SetChat(bot)

	go bot.ProcessMessages()
	go webServer(fakeServer, root)

	fakeServer.Run(context.Background())
}

func webServer(fake *server.FakeSMTP, content fs.FS) { // terrible, but works
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

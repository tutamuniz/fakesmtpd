package main

import (
	"context"
	"embed"
	"io/fs"
	"os"

	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
	"github.com/tutamuniz/fakesmtpd/internal/helper/http"
	"github.com/tutamuniz/fakesmtpd/internal/server"
)

//go:embed web/index.html
var content embed.FS

func main() {
	root, _ := fs.Sub(content, "web")

	path := os.Args[1]

	fakeServer := server.NewServer("0.0.0.0:25", path)

	channel := os.Getenv("CHANNEL_ID")
	apiToken := os.Getenv("API_TOKEN")

	bot := chat.NewBot(apiToken, channel, fakeServer.Logger)

	fakeServer.SetChat(bot)

	go bot.ProcessMessages()
	go http.Server(fakeServer, root)

	fakeServer.Run(context.Background())
}

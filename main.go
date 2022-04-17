package main

import (
	"context"
	"embed"
	"io/fs"
	"os"

	"github.com/tutamuniz/fakesmtpd/helper/chat"
	"github.com/tutamuniz/fakesmtpd/helper/web"
	"github.com/tutamuniz/fakesmtpd/server"
)

//go:embed embfs/index.html
var content embed.FS

func main() {
	root, _ := fs.Sub(content, "embfs")

	path := os.Args[1]

	fakeServer := server.NewServer("0.0.0.0:25", path)

	channel := os.Getenv("CHANNEL_ID")
	api_token := os.Getenv("API_TOKEN")

	bot := chat.NewBot(api_token, channel, fakeServer.Logger)

	fakeServer.SetChat(bot)

	go bot.ProcessMessages()
	go web.Server(fakeServer, root)

	fakeServer.Run(context.Background())
}

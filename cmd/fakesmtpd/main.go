package main

import (
	"context"
	"os"

	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
	"github.com/tutamuniz/fakesmtpd/internal/helper/http"
	"github.com/tutamuniz/fakesmtpd/internal/server"
)

func main() {
	path := os.Args[1]

	fakeServer := server.NewServer("0.0.0.0:25", path)

	channel := os.Getenv("CHANNEL_ID")
	apiToken := os.Getenv("API_TOKEN")

	bot := chat.NewBot(apiToken, channel, fakeServer.Logger)

	fakeServer.SetChat(bot)

	go bot.ProcessMessages()
	go http.Server(fakeServer)

	fakeServer.Run(context.Background())
}

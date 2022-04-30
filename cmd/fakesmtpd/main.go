package main

import (
	"context"
	"flag"
	"log"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
	"github.com/tutamuniz/fakesmtpd/internal/helper/http"
	"github.com/tutamuniz/fakesmtpd/internal/server"
)

var configPath = flag.String("config", "config.toml", "Path to config file")

func main() {
	flag.Parse()

	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	fakeServer := server.NewServer(config)

	bot := chat.NewBot(config.ChatConfig, fakeServer.Logger)

	fakeServer.SetChat(bot)

	go bot.ProcessMessages()
	go http.Server(fakeServer)

	fakeServer.Run(context.Background())
}

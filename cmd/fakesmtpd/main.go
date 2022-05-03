package main

import (
	"flag"
	"log"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/server"
	"github.com/tutamuniz/fakesmtpd/pkg/logging"
)

var configPath = flag.String("config", "config.toml", "Path to config file")

func main() {
	flag.Parse()

	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.AddLogger(logging.NewLogrusLogging())

	fakeServer := server.NewServer(config)

	fakeServer.Run()
}

package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/helper"
	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
	"github.com/tutamuniz/fakesmtpd/internal/helper/http"
)

type FakeSMTP struct {
	capture   bool
	address   string
	wrtimeout int
	rdtimeout int
	config    *config.Config
	helpers   []helper.Helper
	chat      chat.Chat
}

// NewServer init function
func NewServer(config *config.Config) *FakeSMTP {
	logfile := config.LoggingConfig.File

	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	config.Logger = log.New(f, "fakesmtpd ", log.LstdFlags)

	bot := chat.NewBot(config)
	httpserver := http.NewHTTPServer(config)

	return &FakeSMTP{
		capture:   false,
		address:   config.MailServerConfig.Address,
		wrtimeout: 15,
		rdtimeout: 15,
		config:    config,
		helpers:   []helper.Helper{bot, httpserver},
	}
}

func (fake *FakeSMTP) newConnection(c net.Conn) *Connection {
	return &Connection{
		srv:  fake,
		conn: c,
	}
}

func (fake *FakeSMTP) Run() {
	fmt.Printf("Starting server %s\n", fake.address)
	fake.config.Logger.Printf("Starting server %s\n", fake.address)

	for h := range fake.helpers {
		fake.helpers[h].Start()
	}

	server, err := net.Listen("tcp", fake.address)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fake.config.Logger.Printf("ERR: %s", err)
		}
		c := fake.newConnection(conn)
		go c.Handle()

	}
}

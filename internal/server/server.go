package server

import (
	"fmt"
	"net"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/helper"
	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
	"github.com/tutamuniz/fakesmtpd/internal/helper/http"
)

type FakeSMTP struct {
	capture    bool
	address    string
	config     *config.Config
	Chat       helper.Helper
	HTTPServer helper.Helper
}

// NewServer init function
func NewServer(config *config.Config) *FakeSMTP {
	
	bot := chat.NewBot(config)
	httpserver := http.NewHTTPServer(config)

	return &FakeSMTP{
		capture: false,
		address: config.MailServerConfig.Address,

		config:     config,
		Chat:       bot,
		HTTPServer: httpserver,
	}
}

func (fake *FakeSMTP) newConnection(conn net.Conn) *Connection {
	var cc interface{} = fake.Chat

	c, ok := cc.(chat.Chat)

	if !ok {
		panic("chat.Chat interface is not implemented")
	}

	return &Connection{
		config: fake.config,
		chat:   c,
		conn:   conn,
	}
}

func (fake *FakeSMTP) Run() {
	fmt.Printf("Starting server %s\n", fake.address)
	fake.config.Logger.Printf("Starting server %s\n", fake.address)

	go fake.Chat.Start()
	go fake.HTTPServer.Start()

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

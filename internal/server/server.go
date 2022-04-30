package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
)

type FakeSMTP struct {
	capture   bool
	address   string
	wrtimeout int
	rdtimeout int
	datadir   string
	Logger    *log.Logger
	chat      chat.Chat
}

// NewServer init function
func NewServer(address, datadir string) *FakeSMTP {
	logfile := "fakesmtpd.log" // load from config
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(f, "fakesmtpd ", log.LstdFlags)

	return &FakeSMTP{
		capture:   false,
		address:   address,
		wrtimeout: 15,
		rdtimeout: 15,
		datadir:   datadir,
		Logger:    logger,
	}
}

func (fake *FakeSMTP) SetChat(chat chat.Chat) {
	fake.chat = chat
}

func (fake *FakeSMTP) EnableCapture() {
	fake.capture = true
}

func (fake *FakeSMTP) DisableCapture() {
	fake.capture = false
}

func (fake FakeSMTP) CaptureStatus() bool {
	return fake.capture
}

func (fake *FakeSMTP) newConnection(c net.Conn) *Connection {
	return &Connection{
		srv:  fake,
		conn: c,
	}
}

func (fake *FakeSMTP) Run(_ context.Context) {
	fmt.Printf("Starting server %s\n", fake.address)
	fake.Logger.Printf("Starting server %s\n", fake.address)

	server, err := net.Listen("tcp", fake.address)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fake.Logger.Printf("ERR: %s", err)
		}
		c := fake.newConnection(conn)
		go c.Handle()

	}
}

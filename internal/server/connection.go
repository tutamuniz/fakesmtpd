package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/handler"
	"github.com/tutamuniz/fakesmtpd/internal/helper/chat"
)

type Connection struct {
	config *config.Config
	chat   chat.Chat
	conn   net.Conn
	br     *bufio.Reader
	bw     *bufio.Writer
}

func (c Connection) chooseHandler() handler.Handler {
	// read from config
	if !c.config.CaptureStatus {
		return &handler.Default{}
	}

	return &handler.HikVision{
		Logger:  c.config.Logger,
		Chat:    c.chat,
		DataDir: c.config.MailServerConfig.Datadir, // ugly
	}
}

func (c *Connection) Handle() {
	foundDataCmd := false
	defer c.conn.Close()

	handler := c.chooseHandler()

	c.config.Logger.Printf("Using %s\n", handler)
	c.config.Logger.Printf("Connection from: %s\n", c.conn.RemoteAddr().String())

	buff := bytes.Buffer{}
	data := bytes.Buffer{}

	c.br = bufio.NewReader(c.conn)
	c.bw = bufio.NewWriter(c.conn)

	if c.config.MailServerConfig.ReadTimeout != 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.config.MailServerConfig.ReadTimeout)))
	}

	c.write("220 ESMTP  fakesmtp 0.1b\r\n")
	c.bw.Flush()

	for {
		b, err := c.br.ReadSlice('\n')
		if err != nil {
			c.config.Logger.Printf("ERR(ReadingSlice): %s ", err)
			return
		}

		buff.Write(b)
		if foundDataCmd {
			data.Write(b)
		}

		line := string(b)

		if strings.Contains(line, "QUIT") {
			return
		} else if strings.Contains(line, "EHLO") || strings.Contains(line, "HELO") {
			handler.DoHelo("hello")
			c.write("250 DSN")
		} else if strings.Contains(line, "MAIL FROM") {
			handler.DoMailFrom("mail from")
			c.write("250 2.1.0 Ok")
		} else if strings.Contains(line, "RCPT TO") {
			handler.DoRcptTo("rcpto")
			c.write("250 2.1.0 Ok")
		} else if strings.Contains(line, "DATA") {
			foundDataCmd = true
			c.write("354 Go ahead")
		} else if line[0] == '.' && foundDataCmd {
			handler.DoData(data.Bytes())
			c.write("250 2.0.0 Ok: queued")
		}

	}
}

func (c *Connection) write(s string) {
	if c.config.MailServerConfig.WriteTimeout != 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(c.config.MailServerConfig.WriteTimeout)))
	}

	fmt.Fprintf(c.bw, "%s\r\n", s)
	c.bw.Flush()
}

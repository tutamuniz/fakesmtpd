package test

import (
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tutamuniz/fakesmtpd/internal/config"
	"github.com/tutamuniz/fakesmtpd/internal/server"
)

func sendEmailData(file string) {
}

var (
	conf       *config.Config
	fakeServer *server.FakeSMTP
)

func init() {
	var err error
	conf, err = config.LoadConfig("config.toml")

	conf.EnableCapture()

	if err != nil {
		log.Fatal(err)
		return
	}

	fakeServer = server.NewServer(conf)

	go fakeServer.Run()
}

func TestIntegrationSendingFile(t *testing.T) {
	//	t.Parallel()
	port := "25"
	host := "localhost"
	buffer := make([]byte, 1024)

	requests := []struct {
		label  string
		expect string
		cmd    string
	}{
		{"HELO", "220", "HELO localhost"},
		{"MAIL FROM", "250", "MAIL FROM: stranguy@example.com"},
		{"RCPT TO", "250", "RCPT TO: stranguy@example.com"},
		{"DATA", "250", "DATA"},
		{"Send", "354", "Send"},
		{"QUIT", "250", "QUIT"},
	}

	t.Run("request", func(t *testing.T) {
		// Send email
		conn, err := net.Dial("tcp", host+":"+port)
		if err != nil {
			t.Fatal(err)
		}
		idx := 0
		for {

			if idx >= len(requests) {
				break
			}

			nr, _ := conn.Read(buffer)

			if idx != 4 {
				if strings.Contains(string(buffer[:nr]), requests[idx].expect) {
					_, err = conn.Write([]byte(requests[idx].cmd + "\r\n"))

					if err != nil {
						t.Fatal(err)
					}

					idx++
				}
			} else {
				idx++
				content, err := os.ReadFile("email_data_sample.txt")
				if err != nil {
					t.Fatal(err)
				}

				_, err = conn.Write(content)

				if err != nil {
					t.Fatal(err)
				}

			}
			time.Sleep(100 * time.Millisecond)

		}

		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
}

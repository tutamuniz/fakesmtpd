package http

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/tutamuniz/fakesmtpd/internal/config"
)

func TestHttpHelper(t *testing.T) {
	config := &config.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		HTTPServerConfig: config.HTTPServerConfig{
			Address: ":8080",
		},
	}

	httpHelper := NewHTTPServer(config)

	go httpHelper.Start()

	resp, err := http.Get("http://localhost:8080/capture/enable")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	output := string(b)

	if output != "ENABLECAP OK" {
		t.Fatalf("Expected ENABLECAP OK, got %s", output)
	}

	resp, err = http.Get("http://localhost:8080/capture/disable")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	output = string(b)

	if output != "DISABLECAP OK" {
		t.Fatalf("Expected DISABLECAP OK, got %s", output)
	}
}

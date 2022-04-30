package config

import "testing"

func TestConfig(t *testing.T) {
	config := Config{}
	err := config.Load("config_test.toml")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	if config.MailServerConfig.Address != "0.0.0.0:25" {
		t.Errorf("MailServer.Address is not set")
	}

	if config.MailServerConfig.Datadir != "./data" {
		t.Errorf("MailServer.Datadir is not set")
	}

	if config.LoggingConfig.File != "fakesmtpd.log" {
		t.Errorf("Logging.File is not set")
	}

	if config.ChatConfig.ChannelID != "CHANNEL_ID" {
		t.Errorf("Chat.ChannelID is not set")
	}

	if config.ChatConfig.APIToken != "API_TOKEN" {
		t.Errorf("Chat.APIToken is not set")
	}

	if config.HTTPServerConfig.Address != "0.0.0.0:8080" {
		t.Errorf("HTTPServer.Address is not set")
	}
}

package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	MailServerConfig MailServerConfig `toml:"mail_server"`
	LoggingConfig    LoggingConfig    `toml:"logging"`
	ChatConfig       ChatConfig       `toml:"chat"`
	HTTPServerConfig HTTPServerConfig `toml:"http_server"`
	CaptureStatus    bool             `toml:"capture_status"`
	Logger           *log.Logger      `toml:"-"`
}

type MailServerConfig struct {
	Address string
	Datadir string
}

type LoggingConfig struct {
	File string
}

type ChatConfig struct {
	ChannelID string `toml:"channel_id"`
	APIToken  string `toml:"api_token"`
}

type HTTPServerConfig struct {
	Address string
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	err := config.Load(path)
	return config, err
}

func (c *Config) Load(path string) error {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) EnableCapture() {
	c.CaptureStatus = true
}

func (c *Config) DisableCapture() {
	c.CaptureStatus = false
}

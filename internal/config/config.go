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
	Address      string
	Datadir      string
	WriteTimeout int `toml:"write_timeout"`
	ReadTimeout  int `toml:"read_timeout"`
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

	c.validate()

	return nil
}

func (c *Config) EnableCapture() {
	c.CaptureStatus = true
}

func (c *Config) DisableCapture() {
	c.CaptureStatus = false
}

func (c *Config) validate() {
	if c.MailServerConfig.Address == "" {
		c.MailServerConfig.Address = ":25"
	}

	if c.MailServerConfig.WriteTimeout == 0 {
		c.MailServerConfig.WriteTimeout = 15
	}

	if c.MailServerConfig.ReadTimeout == 0 {
		c.MailServerConfig.ReadTimeout = 15
	}

	if c.HTTPServerConfig.Address == "" {
		c.HTTPServerConfig.Address = ":8080"
	}
}

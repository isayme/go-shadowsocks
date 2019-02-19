package conf

import (
	config "github.com/isayme/go-config"
)

// Config config info
type Config struct {
	Method   string `json:"method"`
	Password string `json:"password"`

	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`

	Timeout int `json:"timeout"` // in seconds

	LogLevel string `json:"log_level"`
}

var globalConfig = Config{
	Method:     defaultMethod,
	Server:     defaultServer,
	ServerPort: defaultPort,
	Timeout:    defaultTimeout,
	LogLevel:   defaultLogLevel,
}

// Get parse config
func Get() *Config {
	config.Parse(&globalConfig)
	return &globalConfig
}

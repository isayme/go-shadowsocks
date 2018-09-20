package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// ParseConfig parse config
func ParseConfig(path string) (config *Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config = &Config{
		Method:     defaultMethod,
		Server:     defaultServer,
		ServerPort: defaultPort,
		Timeout:    defaultTimeout,
		LogLevel:   defaultLogLevel,
	}

	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.Password == "" {
		return nil, fmt.Errorf("password required")
	}

	return config, nil
}

package config

import (
	"errors"
	"os"
)

type Config struct {
	ConnString string
}

func GetConfig() (*Config, error) {
	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		return nil, errors.New("CONN_STRING environment variable not set")
	}

	return &Config{connString}, nil
}

package config

import (
	"fmt"
	"io"

	yaml "github.com/goccy/go-yaml"
)

type Config struct {
	Source    Source
	Actions   []Action
	Notifiers []Notifier
}

type Source struct {
	Name    string
	Options map[string]string
}

type Action struct {
	Name    string
	Options map[string]string
}

type Notifier struct {
	Name    string
	Options map[string]string
}

func FromReader(r io.Reader) (Config, error) {
	decoder := yaml.NewDecoder(r)

	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}
	return config, nil
}

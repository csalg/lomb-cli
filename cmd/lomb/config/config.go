package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

func Read() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("getting user config directory: %w", err)
	}
	filename := filepath.Join(configDir, "lomb", "config.json")
	if _, err := os.Stat(filename); err != nil {
		return Config{}, fmt.Errorf("statting config file: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return Config{}, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config file: %w", err)
	}

	return conf, nil
}

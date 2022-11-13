package io

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/csalg/lomb-cli/pkg/types"
)

const (
	configFilename = ".config/lomb/config.json"
)

// ReadConfig reads the configuration file.
func ReadConfig() (config types.Config, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	configPath := filepath.Join(homeDir, configFilename)
	configFile, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	configFileBytes, err := io.ReadAll(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

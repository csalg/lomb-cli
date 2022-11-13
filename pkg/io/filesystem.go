package io

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
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

// ReadLemmas reads a list of lemmas from a file.
func ReadLemmas(filename string) (lemmas []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lemma := scanner.Text()
		if lemma == "" {
			continue
		}
		lemmas = append(lemmas, lemma)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lemmas
}

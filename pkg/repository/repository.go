package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/csalg/lomb-cli/cmd/lomb/cmd/initcmd"
	"github.com/csalg/lomb-cli/pkg/types"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetText(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	return string(b), nil
}

func (r *Repository) SaveAsJSON(filename string, processedText types.ProcessedText) error {
	file, _ := json.MarshalIndent(processedText, "", " ")
	//nolint: gosec
	if err := os.WriteFile(filename, file, 0o644); err != nil {
		return err
	}
	return nil
}

// ReadConfig reads the configuration
func ReadConfig() (types.Config, error) {
	lp, err := ReadLanguagePairConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("reading language pair config: %w", err)
	}
	app, err := ReadAppConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("reading app config: %w", err)
	}
	return types.Config{
		LanguagePairConfig: lp,
		AppConfig:          app,
	}, nil
}

// ReadLanguagePairConfig returns the language pair configuration
func ReadLanguagePairConfig() (*types.LanguagePairConfig, error) {
	return &initcmd.DefaultLanguagePairConfig, nil
}

// ReadAppConfig reads the application configuration
func ReadAppConfig() (types.AppConfig, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return types.AppConfig{}, fmt.Errorf("getting user config directory: %w", err)
	}
	filename := filepath.Join(configDir, "lomb", "config.json")
	if _, err := os.Stat(filename); err != nil {
		return types.AppConfig{}, fmt.Errorf("statting config file: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return types.AppConfig{}, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return types.AppConfig{}, fmt.Errorf("reading config file: %w", err)
	}

	var conf types.AppConfig
	if err := json.Unmarshal(b, &conf); err != nil {
		return types.AppConfig{}, fmt.Errorf("unmarshaling config file: %w", err)
	}
	if conf.Port == "" {
		conf.Port = "9090"
	}
	if conf.BaseURL == "" {
		conf.BaseURL = "http://localhost"
	}

	return conf, nil
}

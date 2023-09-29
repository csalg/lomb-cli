package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

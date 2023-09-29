package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetText(filename string) (string, error) {
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

func SaveAsJSON(filename string, i interface{}) error {
	file, _ := json.MarshalIndent(i, "", " ")
	//nolint: gosec
	if err := os.WriteFile(filename, file, 0o644); err != nil {
		return err
	}
	return nil
}

func ReadAndUnmarshal(filename string, v interface{}) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("checking if config file exists: %w", err)
	}
	if err := unmarshalFile(filename, &v); err != nil {
		return false, fmt.Errorf("unmarshaling file: %w", err)
	}
	return true, nil
}

func unmarshalFile(path_ string, target interface{}) error {
	file, err := os.ReadFile(path_)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path_, err)
	}
	if err := json.Unmarshal(file, &target); err != nil {
		return fmt.Errorf("unmarshalling %s: %w", path_, err)
	}
	return nil
}

package vm

import (
	"os"
	"path/filepath"
)

func getDataDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, ".local", "share", "vm")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func getImagesDir() (string, error) {
	dir, err := getDataDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "images")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return dir, nil
}

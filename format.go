package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var suffixes = map[string]string{
	".raw": "raw",
	".img": "raw",
}

func convertToQcow2(filePath string) error {
	suffix := filepath.Ext(filePath)
	if suffix == "" {
		return fmt.Errorf("%v has no suffix", filePath)
	}

	if suffix == ".qcow2" {
		return nil
	}

	format, ok := suffixes[suffix]
	if !ok {
		return fmt.Errorf("unknown format: %v", suffix)
	}

	qcowFilePath := strings.TrimSuffix(filePath, suffix) + ".qcow2"
	cmd := exec.Command("qemu-img", "convert", "-f", format, "-O", "qcow2", filePath, qcowFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

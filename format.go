package vm

import (
	"os"
	"os/exec"
	"strings"
)

func convertToQcow2(filePath string) error {
	qcowFilePath := strings.TrimSuffix(filePath, ".raw") + ".qcow2"
	cmd := exec.Command("qemu-img", "convert", "-f", "raw", "-O", "qcow2", filePath, qcowFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

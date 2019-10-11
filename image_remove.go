package box

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func imageRemove(c *cli.Context) error {
	name := c.String("name")

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(imagesDir, name+".qcow2")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	if !c.Bool("force") {
		fmt.Printf("Are you sure you want to remove %v? (y/N) ", name+".qcow2")
		var response string
		if _, err := fmt.Scan(&response); err != nil {
			return err
		}
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			return nil
		}
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

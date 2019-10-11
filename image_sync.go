package box

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func imageSync(c *cli.Context) error {
	// TODO: Respect E-Tag
	resp, err := http.Get(baseURL + "index")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(imagesDir, "index"))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return nil
}

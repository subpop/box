package vm

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// ImageGet downloads rawurl and prepares it for use as a backing disk image.
func ImageGet(rawurl string) error {
	URL, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	resp, err := http.Get(URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(imagesDir, filepath.Base(URL.Path))

	if err := download(resp, filePath); err != nil {
		return err
	}

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/x-xz":
		if err := decompressXZ(filePath); err != nil {
			return err
		}
		filePath = strings.TrimSuffix(filePath, ".xz")
	case "application/gzip":
		if err := decompressGZ(filePath); err != nil {
			return err
		}
		filePath = strings.TrimSuffix(filePath, ".gz")
	}

	// TODO: fix this switch to detect raw files after download
	if err := convertToQcow2(filePath); err != nil {
		return err
	}

	return nil
}

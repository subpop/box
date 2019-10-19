package vm

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// TemplateGet downloads and prepares a disk template for use as a backing disk image.
func TemplateGet(name, arch string) error {
	index, err := newIndex()
	if err != nil {
		return err
	}

	var template template
	for _, t := range index.Templates {
		if t.ININame == name && t.Arch == arch {
			template = t
			break
		}
	}

	URL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	URL.Path += template.File

	resp, err := http.Get(URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(imagesDir, template.File)

	if err := download(resp, filePath); err != nil {
		return err
	}

	checksum, ok := template.Checksum["sha512"]
	if !ok {
		checksum = template.Checksum[""]
	}
	if err := verify(filePath, checksum); err != nil {
		return err
	}

	if err := decompressXZ(filePath); err != nil {
		return err
	}

	if err := os.Rename(strings.TrimSuffix(filePath, ".xz"), filepath.Join(imagesDir, strings.TrimSuffix(template.File, ".xz")+"."+template.Format)); err != nil {
		return err
	}

	filePath = filepath.Join(imagesDir, strings.TrimSuffix(template.File, ".xz")+"."+template.Format)

	if err := convertToQcow2(filePath); err != nil {
		return err
	}

	return nil
}

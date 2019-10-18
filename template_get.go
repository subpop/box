package vm

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uiprogress"
	"github.com/ulikunitz/xz"
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

	contentLengthHeader := resp.Header.Get("Content-Length")
	if contentLengthHeader == "" {
		return errors.New("missing header: Content-Length")
	}

	contentLength, err := strconv.ParseInt(contentLengthHeader, 10, 64)
	if err != nil {
		return err
	}

	uiprogress.Start()
	bar := uiprogress.AddBar(int(contentLength))
	bar.AppendCompleted()

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(imagesDir, template.File)

	f, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	bc := &byteCounter{
		bar: bar,
	}
	defer bc.Close()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("dl: %s", humanize.Bytes(bc.total))
	})

	_, err = io.Copy(f, io.TeeReader(resp.Body, bc))
	if err != nil {
		return err
	}

	err = os.Rename(filePath+".tmp", filePath)
	if err != nil {
		return err
	}

	r, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	hash := sha512.New()
	_, err = hash.Write(data)
	if err != nil {
		return err
	}
	computed := fmt.Sprintf("%x", hash.Sum(nil))
	checksum, ok := template.Checksum["sha512"]
	if !ok {
		checksum = template.Checksum[""]
	}
	if checksum != computed {
		return fmt.Errorf("invalid checksum: %v != %v", checksum, computed)
	}

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/x-xz":
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer os.Remove(filePath)

		filePath = strings.TrimSuffix(filePath, ".xz") + "." + template.Format
		tmpFilePath := filePath + ".tmp"
		xzf, err := os.Create(tmpFilePath)
		if err != nil {
			return err
		}
		defer xzf.Close()

		xzr, err := xz.NewReader(f)
		if err != nil {
			return err
		}

		bar := uiprogress.AddBar(int(template.Size))
		bar.AppendCompleted()

		bc := &byteCounter{
			bar: bar,
		}
		defer bc.Close()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("xz: %s", humanize.Bytes(bc.total))
		})

		_, err = io.Copy(xzf, io.TeeReader(xzr, bc))
		if err != nil {
			return err
		}

		err = os.Rename(tmpFilePath, filePath)
		if err != nil {
			return err
		}
	}

	switch template.Format {
	case "raw":
		qcowFilePath := strings.TrimSuffix(filePath, ".raw") + ".qcow2"
		cmd := exec.Command("qemu-img", "convert", "-f", "raw", "-O", "qcow2", filePath, qcowFilePath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	uiprogress.Stop()

	return nil
}

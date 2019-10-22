package vm

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/ulikunitz/xz"
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

func getInstancesDir() (string, error) {
	dir, err := getDataDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "instances")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func inspect(filePath string) error {
	switch filepath.Ext(filePath) {
	case ".gz", ".xz":
		return decompress(filePath)
	case ".raw", ".img":
		return convert(filePath)
	case ".qcow2":
		return nil
	}

	return fmt.Errorf("unsupported file type: %v", filePath)
}

func download(rawurl string) (string, error) {
	URL, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(URL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imagesDir, err := getImagesDir()
	if err != nil {
		return "", err
	}

	destFilePath := filepath.Join(imagesDir, filepath.Base(URL.Path))
	w, err := os.Create(destFilePath + ".tmp")
	if err != nil {
		return "", err
	}
	defer w.Close()

	var bytesWritten uint64
	err = copy(w, resp.Body, func(buf []byte) {
		bytesWritten += uint64(len(buf))
		fmt.Printf("\r%s", strings.Repeat(" ", 40))
		fmt.Printf("\rdownloading... %s", humanize.Bytes(bytesWritten))
	})
	if err != nil {
		return "", err
	}

	err = os.Rename(destFilePath+".tmp", destFilePath)
	if err != nil {
		return "", err
	}

	return destFilePath, nil
}

func decompress(filePath string) error {
	var err error

	r, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer r.Close()

	var g io.Reader
	switch filepath.Ext(filePath) {
	case ".gz":
		g, err = gzip.NewReader(r)
	case ".xz":
		g, err = xz.NewReader(r)
	}
	if err != nil {
		return err
	}

	destFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	w, err := os.Create(destFilePath + ".tmp")
	if err != nil {
		return err
	}
	defer w.Close()

	var bytesWritten uint64
	err = copy(w, g, func(buf []byte) {
		bytesWritten += uint64(len(buf))
		fmt.Printf("\r%s", strings.Repeat(" ", 40))
		fmt.Printf("\rdecompressing... %s", humanize.Bytes(bytesWritten))
	})
	if err != nil {
		return err
	}

	if err := os.Rename(destFilePath+".tmp", destFilePath); err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return inspect(destFilePath)
}

func convert(filePath string) error {
	var err error

	destFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".qcow2"
	cmd := exec.Command("qemu-img", "convert",
		"-f", "raw",
		"-O", "qcow2",
		filePath, destFilePath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return inspect(destFilePath)
}

func copy(dest io.Writer, src io.Reader, writeFunc func(buf []byte)) error {
	w := printWriter{
		print: writeFunc,
	}

	_, err := io.Copy(dest, io.TeeReader(src, w))
	if err != nil {
		return err
	}
	return nil
}

type printWriter struct {
	print func(buf []byte)
}

func (p printWriter) Write(buf []byte) (int, error) {
	p.print(buf)
	return len(buf), nil
}

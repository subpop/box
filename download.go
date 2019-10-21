package vm

import (
	"compress/gzip"
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/ulikunitz/xz"
)

func download(resp *http.Response, filePath string) error {
	var err error

	f, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	var bytesWritten uint64
	bc := &byteCounter{
		onWrite: func(buf []byte) {
			bytesWritten += uint64(len(buf))
			fmt.Printf("\rdownloading... %s", humanize.Bytes(bytesWritten))
		},
		onClose: func() {
			fmt.Println()
		},
	}

	_, err = io.Copy(f, io.TeeReader(resp.Body, bc))
	if err != nil {
		return nil
	}
	bc.Close()

	err = os.Rename(filePath+".tmp", filePath)
	if err != nil {
		return err
	}

	return nil
}

func verify(filePath, checksum string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	hash := sha512.New()
	_, err = hash.Write(data)
	if err != nil {
		return err
	}
	computed := fmt.Sprintf("%x", hash.Sum(nil))
	if checksum != computed {
		return fmt.Errorf("invalid checksum: %v != %v", checksum, computed)
	}

	return nil
}

func decompressGZ(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer os.Remove(filePath)

	filePath = strings.TrimSuffix(filePath, ".gz")
	tmpFilePath := filePath + ".tmp"
	gzf, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer gzf.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	if err := copy(gzf, gzr); err != nil {
		return err
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}

	return nil
}

func decompressXZ(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer os.Remove(filePath)

	filePath = strings.TrimSuffix(filePath, ".xz")
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

	if err := copy(xzf, xzr); err != nil {
		return err
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}

	return nil
}

func copy(dest io.Writer, src io.Reader) error {
	var err error

	var bytesWritten uint64
	bc := &byteCounter{
		onWrite: func(buf []byte) {
			bytesWritten += uint64(len(buf))
			fmt.Printf("\rdecompressing... %s", humanize.Bytes(bytesWritten))
		},
		onClose: func() {
			fmt.Println()
		},
	}
	defer bc.Close()

	_, err = io.Copy(dest, io.TeeReader(src, bc))
	if err != nil {
		return err
	}

	return nil
}

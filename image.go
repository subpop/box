package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/subpop/go-ini"
)

const baseURL string = "http://builder.libguestfs.org/"

type image struct {
	ININame        string
	Name           string            `ini:"name"`
	OSInfo         string            `ini:"osinfo,omitempty"`
	Arch           string            `ini:"arch"`
	File           string            `ini:"file"`
	Revision       int               `ini:"revision,omitempty"`
	Checksum       map[string]string `ini:"checksum"`
	Format         string            `ini:"format"`
	Size           uint64            `ini:"size"`
	CompressedSize uint64            `ini:"compressed_size"`
	Expand         string            `ini:"expand"`
	Notes          string            `ini:"notes"`
}

type index struct {
	Images []image `ini:"*"`
}

func newIndex() (i index, err error) {
	imagesDir, err := getImagesDir()
	if err != nil {
		return
	}

	f, err := os.Open(filepath.Join(imagesDir, "index"))
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = ini.UnmarshalWithOptions(data, &i, ini.Options{AllowMultilineValues: true})
	if err != nil {
		return
	}

	return
}

// getImage looks up an image by name and arch in the index
func (i *index) getImage(name, arch string) (*image, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if arch == "" {
		arch = "x86_64"
	}

	for _, image := range i.Images {
		if image.ININame == name && image.Arch == arch {
			return &image, nil
		}
	}

	return nil, fmt.Errorf("no image with name %q and arch %q", name, arch)
}

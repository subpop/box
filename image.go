package box

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

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

func (i image) String() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\t%v\n", i.Name)
	fmt.Fprintf(w, "OSINFO\t%v\n", i.OSInfo)
	fmt.Fprintf(w, "ARCH\t%v\n", i.Arch)
	fmt.Fprintf(w, "FILE\t%v\n", i.File)
	fmt.Fprintf(w, "REVISION\t%v\n", i.Revision)
	fmt.Fprintf(w, "CHECKSUM\t%v\n", i.Checksum)
	fmt.Fprintf(w, "FORMAT\t%v\n", i.Format)
	fmt.Fprintf(w, "SIZE\t%v\n", i.Size)
	fmt.Fprintf(w, "NOTES\n%v\n", i.Notes)
	w.Flush()
	return b.String()
}

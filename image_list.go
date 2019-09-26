package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/subpop/go-ini"

	"github.com/urfave/cli"
)

type image struct {
	SectionName    string
	Name           string `ini:"name"`
	OSInfo         string `ini:"osinfo,omitempty"`
	Arch           string `ini:"arch"`
	File           string `ini:"file"`
	Revision       int    `ini:"revision,omitempty"`
	Checksum       string `ini:"checksum,omitempty"`
	Format         string `ini:"format"`
	Size           uint64 `ini:"size"`
	CompressedSize uint64 `ini:"compressed_size"`
	Expand         string `ini:"expand"`
}

func imageList(c *cli.Context) (err error) {
	res, err := http.Get("http://builder.libguestfs.org/index")
	if err != nil {
		return
	}
	defer res.Body.Close()

	var imageIndex struct {
		Images []image `ini:"[.*]"`
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = ini.UnmarshalWithOptions(data, &imageIndex, ini.Options{AllowMultilineValues: true})
	if err != nil {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tARCH\tDESCRIPTION\t")
	for _, image := range imageIndex.Images {
		fmt.Fprintf(w, "%v\t%v\t%v\n", image.SectionName, image.Arch, image.Name)
	}
	w.Flush()

	return
}

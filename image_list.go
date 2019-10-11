package box

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

func imageList(c *cli.Context) error {
	index, err := newIndex()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tARCH\tDESCRIPTION\t")
	for _, image := range index.Images {
		fmt.Fprintf(w, "%v\t%v\t%v\n", image.ININame, image.Arch, image.Name)
	}
	w.Flush()

	return nil
}

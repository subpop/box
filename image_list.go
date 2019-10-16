package vm

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// ImageList prints a list of available images.
func ImageList(sortBy string) error {
	index, err := newIndex()
	if err != nil {
		return err
	}

	sort.Slice(index.Images, func(i, j int) bool {
		switch sortBy {
		case "arch":
			return index.Images[i].Arch < index.Images[j].Arch
		case "description", "desc":
			return index.Images[i].Name < index.Images[j].Name
		case "name":
			fallthrough
		default:
			return index.Images[i].ININame < index.Images[j].ININame
		}
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tARCH\tDESCRIPTION\t")
	for _, image := range index.Images {
		fmt.Fprintf(w, "%v\t%v\t%v\n", image.ININame, image.Arch, image.Name)
	}
	w.Flush()

	return nil
}

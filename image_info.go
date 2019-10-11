package box

import (
	"fmt"

	"github.com/urfave/cli"
)

func imageInfo(c *cli.Context) error {
	index, err := newIndex()
	if err != nil {
		return err
	}

	name := c.String("name")
	arch := c.String("arch")

	var image image
	for _, i := range index.Images {
		if i.ININame == name && i.Arch == arch {
			image = i
			break
		}
	}

	fmt.Println(image)

	return nil
}

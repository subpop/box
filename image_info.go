package vm

import (
	"fmt"
)

// ImageInfo prints information about image name.
func ImageInfo(name, arch string) error {
	index, err := newIndex()
	if err != nil {
		return err
	}

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

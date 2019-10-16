package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/libvirt/libvirt-go"
)

// Destroy stops and undefines a domain by name. If force is true, the
// domain is destroyed without prompting for confirmation.
func Destroy(name string, force bool) error {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}

	dom, err := conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer dom.Free()

	name, err = dom.GetName()
	if err != nil {
		return err
	}

	if !force {
		fmt.Printf("Are you sure you wish to destroy %v? (y/N) ", name)
		var response string
		fmt.Scan(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" {
			return nil
		}
	}

	state, _, err := dom.GetState()
	if err != nil {
		return err
	}
	if state == libvirt.DOMAIN_RUNNING {
		err = dom.Destroy()
		if err != nil {
			return err
		}
	}

	UUID, err := dom.GetUUIDString()
	if err != nil {
		return err
	}

	os.Remove(filepath.Join(imagesDir, UUID+".qcow2"))

	err = dom.Undefine()
	if err != nil {
		return err
	}

	return nil
}

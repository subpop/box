package vm

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
)

// Down looks up an active domain by name and stops it. If force is true, the
// domain is stopped ungracefully.
func Down(name string, id int, force bool) error {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	if len(name) > 0 && id > 0 {
		return fmt.Errorf("conflicting arguments: name, id")
	}

	var dom *libvirt.Domain
	if len(name) > 0 {
		dom, err = conn.LookupDomainByName(name)
		if err != nil {
			return nil
		}
	} else if id > 0 {
		dom, err = conn.LookupDomainById(uint32(id))
		if err != nil {
			return nil
		}
	} else {
		return fmt.Errorf("conflicting arguments: name, id")
	}
	defer dom.Free()

	if force {
		if err := dom.Destroy(); err != nil {
			return err
		}
	} else {
		if err := dom.Shutdown(); err != nil {
			return err
		}
	}

	return nil
}

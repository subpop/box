package box

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
	"github.com/urfave/cli"
)

func up(c *cli.Context) error {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	dom, err := conn.LookupDomainByName(c.String("name"))
	if err != nil {
		return err
	}
	defer dom.Free()

	state, _, err := dom.GetState()
	if err != nil {
		return err
	}

	switch state {
	case libvirt.DOMAIN_SHUTOFF:
		if err := dom.Create(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("error: cannot start box in state: %v", state)
	}

	return nil
}

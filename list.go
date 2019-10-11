package box

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/libvirt/libvirt-go"
)

// List prints a list of known domains. If active is true, active domains are
// included in the list. If inactive is true, inactive domains are included
// in the list.
func List(active, inactive bool) error {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}
	defer conn.Close()

	var flags libvirt.ConnectListAllDomainsFlags
	if active {
		flags |= libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	}
	if inactive {
		flags |= libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	}

	domains, err := conn.ListAllDomains(flags)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\t")
	for _, domain := range domains {
		defer domain.Free()
		name, err := domain.GetName()
		if err != nil {
			return err
		}
		id := func() string {
			state, _, err := domain.GetState()
			if err != nil {
				return "-"
			}
			if state == libvirt.DOMAIN_RUNNING {
				id, err := domain.GetID()
				if err != nil {
					return "-"
				}
				return strconv.FormatUint(uint64(id), 10)
			}
			return "-"
		}()
		fmt.Fprintf(w, "%v\t%v\t\n", id, name)
	}
	w.Flush()

	return nil
}

package vm

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/libvirt/libvirt-go"
)

// Info prints detailed information about the given domain.
func Info(name string) error {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	dom, err := conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer dom.Free()

	data, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		return err
	}

	var d domain
	err = xml.Unmarshal([]byte(data), &d)
	if err != nil {
		return err
	}

	fmt.Println(strings.TrimSpace(d.String()))

	return nil
}

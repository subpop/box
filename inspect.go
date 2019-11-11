package vm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/libvirt/libvirt-go"
)

// Inspect prints detailed information about the given domain.
func Inspect(name, outputformat string) error {
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

	var output []byte
	switch outputformat {
	case "xml":
		output, err = xml.MarshalIndent(d, "", "\t")
		if err != nil {
			return err
		}
	case "json":
		output, err = json.MarshalIndent(d, "", "\t")
		if err != nil {
			return err
		}
	}
	fmt.Println(string(output))

	return nil
}

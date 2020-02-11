package vm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/libvirt/libvirt-go"
)

// DomainCapabilities prints detailed information about the domain capabilities.
func DomainCapabilities(outputformat string) error {
	cap, err := getDomainCapabilities()
	if err != nil {
		return err
	}

	var output []byte
	switch outputformat {
	case "json":
		output, err = json.MarshalIndent(cap, "", "\t")
		if err != nil {
			return err
		}
	case "xml":
		output, err = xml.MarshalIndent(cap, "", "\t")
		if err != nil {
			return err
		}
	default:
		return UnsupportedFormatErr{outputformat}
	}
	fmt.Println(strings.TrimSpace(string(output)))

	return nil
}

func getDomainCapabilities() (*domainCapabilities, error) {
	conn, err := libvirt.NewConnect("")
	if err != nil {
		return nil, err
	}

	capabilities, err := getCapabilities()
	if err != nil {
		return nil, err
	}

	arch := capabilities.Host.CPU.Arch
	var emulator string
	var machine string
	for _, guest := range capabilities.Guest {
		if guest.Arch.Name == arch {
			emulator = guest.Arch.Emulator
			for _, m := range guest.Arch.Machine {
				if m.CharData == "pc" {
					machine = m.Canonical
					break
				}
			}
			break
		}
	}

	data, err := conn.GetDomainCapabilities(emulator, arch, machine, "kvm", 0)
	if err != nil {
		return nil, err
	}

	var domcap domainCapabilities
	err = xml.Unmarshal([]byte(data), &domcap)
	if err != nil {
		return nil, err
	}

	return &domcap, nil
}

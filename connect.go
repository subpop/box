package main

import (
	"fmt"
	"io"
	"os"

	"github.com/libvirt/libvirt-go"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func connect(c *cli.Context) error {
	var err error

	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	var dom *libvirt.Domain
	dom, err = conn.LookupDomainByName(c.String("name"))
	if err != nil {
		return err
	}
	defer dom.Free()

	mode := c.String("mode")
	switch mode {
	case "ssh":
		return connectSSH(c, dom)
	case "console":
		return connectConsole(c, dom)
	default:
		return fmt.Errorf("error: unsupported connection mode: %v", mode)
	}
}

func connectSSH(c *cli.Context, dom *libvirt.Domain) error {
	config := &ssh.ClientConfig{
		User: c.String("user"),
		Auth: []ssh.AuthMethod{
			// TODO: Add PublicKeyAuthentication
			ssh.RetryableAuthMethod(ssh.PasswordCallback(func() (secret string, err error) {
				fmt.Print("Password: ")
				data, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return "", err
				}
				return string(data), nil
			}), 3),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	interfaces, err := dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP | libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
	if err != nil {
		return err
	}
	if len(interfaces) == 0 {
		name, err := dom.GetName()
		if err != nil {
			return err
		}
		return fmt.Errorf("error: no interfaces detected for %v", name)
	}
	addrs := make([]string, 0)
	for _, iface := range interfaces {
		for _, addr := range iface.Addrs {
			if addr.Type == int(libvirt.IP_ADDR_TYPE_IPV4) {
				addrs = append(addrs, addr.Addr)
			}
		}
	}

	var addr string
	if len(addrs) > 1 {
		fmt.Println("Multiple addresses detected.")
		for i, addr := range addrs {
			fmt.Printf("%v: %v\n", i, addr)
		}
		fmt.Println("Select address: ")
		var response int
		fmt.Scan(&response)
		addr = addrs[response-1]
	} else {
		addr = addrs[0]
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", addr), config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, stderr)

	fd := int(os.Stdin.Fd())
	if terminal.IsTerminal(fd) {
		oldState, err := terminal.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer terminal.Restore(fd, oldState)

		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			return err
		}

		if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
			return err
		}
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func connectConsole(c *cli.Context, dom *libvirt.Domain) error {
	return nil
}

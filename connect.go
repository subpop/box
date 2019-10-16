package box

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/libvirt/libvirt-go"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const escapeSequence = byte(']') ^ 0x40

// Connect opens a connection to a domain by name. The mode argument determines
// the connection mode: either "ssh" or "console".
func Connect(name string, mode string, user string) error {
	var err error

	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	var dom *libvirt.Domain
	dom, err = conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer dom.Free()

	switch mode {
	case "ssh":
		return connectSSH(dom, user)
	case "console":
		return connectConsole(dom)
	default:
		return fmt.Errorf("error: unsupported connection mode: %v", mode)
	}
}

func connectSSH(dom *libvirt.Domain, user string) error {
	config := &ssh.ClientConfig{
		User: user,
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
			fmt.Printf("%v: %v\n", i+1, addr)
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

func connectConsole(dom *libvirt.Domain) error {
	var err error

	name, err := dom.GetName()
	if err != nil {
		return err
	}
	fmt.Println("Connected to " + name)
	fmt.Println("Escape character is ^]")

	conn, err := dom.DomainGetConnect()
	if err != nil {
		return err
	}

	oldstate, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldstate)

	signal.Ignore(syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGPIPE)

	stream, err := conn.NewStream(0)
	if err != nil {
		return err
	}

	if err := dom.OpenConsole("", stream, libvirt.DOMAIN_CONSOLE_SAFE); err != nil {
		return err
	}

	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	defer cond.L.Unlock()
	var quit bool

	stdin := bufio.NewReader(os.Stdin)

	// read from stream and write to stdout
	go func() {
		var err error
		for !quit {
			var buf []byte
			var got, sent int

			buf = make([]byte, 1024)

			// read from the stream, continuing if no bytes are read
			got, err = stream.Recv(buf)
			if got == 0 {
				if err != nil {
					break
				}
				continue
			}

			// write to stdout
			sent, err = os.Stdout.Write(buf)
			if sent != len(buf) {
				if err != nil {
					break
				}
			}
		}
		if err != nil {
			fmt.Println(err)
		}
		quit = true
		cond.Broadcast()
	}()

	// read from stdin and write to stream
	go func() {
		var err error
		for !quit {
			var got byte

			got, err = stdin.ReadByte()
			if err != nil {
				break
			}

			if got == escapeSequence {
				break
			}

			_, err = stream.Send([]byte{got})
			if err != nil {
				break
			}
		}
		if err != nil {
			fmt.Println(err)
		}
		quit = true
		cond.Broadcast()
	}()

	for !quit {
		cond.Wait()
	}

	signal.Reset()

	return nil
}

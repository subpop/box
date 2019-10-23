package vm

import (
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Create defines a new domain using name and creating a disk image backed by
// image.
func Create(name, image string, disks []string) error {
	if name == "" {
		name = petname.Generate(2, "-")
	}

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}

	instancesDir, err := getInstancesDir()
	if err != nil {
		return err
	}

	baseImagePath := filepath.Join(imagesDir, image+".qcow2")
	if _, err := os.Stat(baseImagePath); os.IsNotExist(err) {
		return err
	}

	var domain domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		return err
	}
	domain.UUID = uuid.New().String()
	domain.Name = name

	overlayImagePath := filepath.Join(instancesDir, domain.UUID+".qcow2")
	cmd := exec.Command("qemu-img",
		"create",
		"-f",
		"qcow2",
		"-o",
		fmt.Sprintf("backing_file=%v", baseImagePath),
		overlayImagePath)
	if err := cmd.Run(); err != nil {
		return err
	}
	domain.Devices.Disks[0].Source.File = overlayImagePath

	instanceDataDir := filepath.Join(instancesDir, domain.UUID)
	if _, err := os.Stat(instanceDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(instanceDataDir, 0755); err != nil {
			return err
		}
	}

	for _, d := range disks {
		src, err := os.Open(d)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.Create(filepath.Join(instanceDataDir, filepath.Base(d)))
		if err != nil {
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			return err
		}

		var device string
		switch filepath.Ext(d) {
		case ".iso":
			device = "cdrom"
		case ".img":
			device = "floppy"
		default:
			device = "disk"
		}

		var drv driver
		switch filepath.Ext(d) {
		case ".qcow2":
			drv = driver{
				Name: "qemu",
				Type: "qcow2",
			}
		default:
			drv = driver{
				Name: "qemu",
				Type: "raw",
			}
		}

		disk := disk{
			Type:   "file",
			Device: device,
			Driver: drv,
			Source: source{
				File: dest.Name(),
			},
			Target: target{
				Dev: "hdb",
				Bus: "ide",
			},
		}
		domain.Devices.Disks = append(domain.Devices.Disks, disk)
	}

	data, err := xml.Marshal(domain)
	if err != nil {
		return err
	}

	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}
	defer conn.Close()

	dom, err := conn.DomainDefineXML(string(data))
	if err != nil {
		return err
	}
	defer dom.Free()

	if err := dom.Create(); err != nil {
		return err
	}

	fmt.Println("Created " + domain.Name)

	return nil
}

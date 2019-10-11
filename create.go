package box

import (
	"encoding/xml"
	"fmt"
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

// Create defines and starts a new box with name and base image.
func Create(name, image string) error {
	if name == "" {
		name = petname.Generate(2, "-")
	}

	imagesDir, err := getImagesDir()
	if err != nil {
		return err
	}

	baseImagePath := filepath.Join(imagesDir, image+".qcow2")
	if _, err := os.Stat(baseImagePath); os.IsNotExist(err) {
		return err
	}

	var domain Domain
	if err := xml.Unmarshal([]byte(domainXML), &domain); err != nil {
		return err
	}
	domain.UUID = uuid.New().String()
	domain.Name = name

	overlayImagePath := filepath.Join(imagesDir, domain.UUID+".qcow2")
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

	return nil
}

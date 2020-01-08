package vm

import (
	"encoding/xml"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
)

type videoAcceleration struct {
	Accel3d string `xml:"accel3d,attr"`
}

type videoModel struct {
	Type         string            `xml:"type,attr,omitempty"`
	VRAM         string            `xml:"vram,attr,omitempty"`
	Heads        string            `xml:"heads,attr,omitempty"`
	Acceleration videoAcceleration `xml:"acceleration,omitempty"`
}

type video struct {
	Model videoModel `xml:"model,omitempty"`
}

type consoleTarget struct {
	Type  string `xml:"type,attr"`
	Alias string `xml:"alias"`
}

type console struct {
	Type   string        `xml:"type,attr"`
	Target consoleTarget `xml:"target"`
}

type netModel struct {
	Type string `xml:"type,attr"`
}

type mac struct {
	Address string `xml:"address,attr"`
}

type interfaceSource struct {
	Network string `xml:"network,attr,omitempty"`
	Bridge  string `xml:"bridge,attr,omitempty"`
}

type netInterface struct {
	Type   string          `xml:"type,attr"`
	Source interfaceSource `xml:"source"`
	MAC    *mac            `xml:"mac,omitempty" json:",omitempty"`
	Model  netModel        `xml:"model"`
}

type master struct {
	StartPort string `xml:"startport,attr,omitempty"`
}

type controller struct {
	Type   string  `xml:"type,attr"`
	Index  string  `xml:"index,attr"`
	Model  string  `xml:"model,attr"`
	Master *master `xml:"master,omitempty" json:",omitempty"`
}

type target struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type source struct {
	File string `xml:"file,attr"`
}

type driver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type disk struct {
	Type     string `xml:"type,attr"`
	Device   string `xml:"device,attr"`
	Driver   driver `xml:"driver"`
	Source   source `xml:"source"`
	Target   target `xml:"target"`
	ReadOnly string `xml:"readonly,omitempty"`
}

type devices struct {
	Emulator    string         `xml:"emulator"`
	Disks       []disk         `xml:"disk"`
	Controllers []controller   `xml:"controller"`
	Interfaces  []netInterface `xml:"interface"`
	Consoles    []console      `xml:"console"`
	Videos      []video        `xml:"video"`
}

type suspendTo struct {
	Enabled string `xml:"enabled,attr"`
}

type pm struct {
	SuspendToMem  suspendTo `xml:"suspend-to-mem"`
	SuspendToDisk suspendTo `xml:"suspend-to-disk"`
}

type timer struct {
	Name       string `xml:"name,attr"`
	TickPolicy string `xml:"tickpolicy,attr,omitempty"`
	Present    string `xml:"present,attr,omitempty"`
}

type clock struct {
	Offset string  `xml:"offset,attr"`
	Timers []timer `xml:"timer"`
}

type cpu struct {
	Mode string `xml:"mode,attr"`
}

type features struct {
	Acpi string `xml:"acpi"`
	Apic string `xml:"apic"`
}

type loader struct {
	ReadOnly string `xml:"readonly,attr"`
	Type     string `xml:"type,attr"`
	Value    string `xml:",chardata"`
}

type boot struct {
	Dev string `xml:"dev,attr"`
}

type osType struct {
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
	Value   string `xml:",chardata"`
}

type operatingSystem struct {
	Type     osType  `xml:"type"`
	Boot     boot    `xml:"boot"`
	Firmware string  `xml:"firmware,attr,omitempty"`
	Loader   *loader `xml:"loader,omitempty" json:",omitempty"`
}

type domain struct {
	XMLName       xml.Name        `xml:"domain" json:"-"`
	Type          string          `xml:"type,attr"`
	Name          string          `xml:"name"`
	UUID          string          `xml:"uuid"`
	Memory        uint            `xml:"memory"`
	CurrentMemory uint            `xml:"currentMemory"`
	VCPU          uint            `xml:"vcpu"`
	OS            operatingSystem `xml:"os"`
	Features      features        `xml:"features"`
	CPU           cpu             `xml:"cpu"`
	Clock         clock           `xml:"clock"`
	PM            pm              `xml:"pm"`
	Devices       devices         `xml:"devices"`
}

func (d domain) String() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\t%v\n", d.Name)
	fmt.Fprintf(w, "ARCH\t%v\n", d.OS.Type.Arch)
	fmt.Fprintf(w, "VCPU\t%v\n", d.VCPU)
	fmt.Fprintf(w, "MEMORY\t%v\n", humanize.Bytes(uint64(d.Memory)))
	for i, disk := range d.Devices.Disks {
		fmt.Fprintf(w, "DISK%v\t%v\n", i, disk.Source.File)
	}
	w.Flush()
	return b.String()
}

const domainXML string = `
<domain type="kvm">
  <name></name>
  <uuid></uuid>
  <memory>524288</memory>
  <currentMemory>524288</currentMemory>
  <vcpu>1</vcpu>
  <os>
    <type arch="x86_64" machine="pc">hvm</type>
    <boot dev="hd"/>
  </os>
  <features>
    <acpi/>
    <apic/>
  </features>
  <cpu mode="host-model"/>
  <clock offset="utc">
    <timer name="rtc" tickpolicy="catchup"/>
    <timer name="pit" tickpolicy="delay"/>
    <timer name="hpet" present="no"/>
  </clock>
  <pm>
    <suspend-to-mem enabled="no"/>
    <suspend-to-disk enabled="no"/>
  </pm>
  <devices>
    <emulator>/usr/bin/qemu-kvm</emulator>
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2"/>
      <source file=""/>
      <target dev="hda" bus="ide"/>
    </disk>
    <controller type="usb" index="0" model="ich9-ehci1"/>
    <controller type="usb" index="0" model="ich9-uhci1">
      <master startport="0"/>
    </controller>
    <controller type="usb" index="0" model="ich9-uhci2">
      <master startport="2"/>
    </controller>
    <controller type="usb" index="0" model="ich9-uhci3">
      <master startport="4"/>
    </controller>
    <interface type="bridge">
      <source bridge="virbr0"/>
      <model type="virtio"/>
    </interface>
	<console type="pty">
	  <target type="serial"/>
	</console>
	<console type="pty">
	  <target type="virtio"/>
	</console>
	<video/>
  </devices>
</domain>`

type domainSnapshot struct {
	XMLName      xml.Name `xml:"domainsnapshot"`
	Name         string   `xml:"name"`
	Description  string   `xml:"description"`
	CreationTime string   `xml:"creationTime"`
	State        string   `xml:"state"`
}

const domainSnapshotXML string = `
<domainsnapshot>
  <name/>
  <description/>
  <disks/>
</domainsnapshot>`

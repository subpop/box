package vm

import (
	"encoding/xml"
)

type Console struct {
	Type string `xml:"type,attr"`
}

type Model struct {
	Type string `xml:"type,attr"`
}

type Mac struct {
	Address string `xml:"address,attr"`
}

type InterfaceSource struct {
	Network string `xml:"network,attr,omitempty"`
	Bridge  string `xml:"bridge,attr,omitempty"`
}

type Interface struct {
	Type   string          `xml:"type,attr"`
	Source InterfaceSource `xml:"source"`
	MAC    *Mac            `xml:"mac,omitempty"`
	Model  Model           `xml:"model"`
}

type Master struct {
	StartPort string `xml:"startport,attr,omitempty"`
}

type Controller struct {
	Type   string  `xml:"type,attr"`
	Index  string  `xml:"index,attr"`
	Model  string  `xml:"model,attr"`
	Master *Master `xml:"master,omitempty"`
}

type Target struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type Source struct {
	File string `xml:"file,attr"`
}

type Driver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Disk struct {
	Type     string `xml:"type,attr"`
	Device   string `xml:"device,attr"`
	Driver   Driver `xml:"driver"`
	Source   Source `xml:"source"`
	Target   Target `xml:"target"`
	ReadOnly string `xml:"readonly,omitempty"`
}

type Devices struct {
	Emulator    string       `xml:"emulator"`
	Disks       []Disk       `xml:"disk"`
	Controllers []Controller `xml:"controller"`
	Interfaces  []Interface  `xml:"interface"`
	Consoles    []Console    `xml:"console"`
}

type SuspendTo struct {
	Enabled string `xml:"enabled,attr"`
}

type Pm struct {
	SuspendToMem  SuspendTo `xml:"suspend-to-mem"`
	SuspendToDisk SuspendTo `xml:"suspend-to-disk"`
}

type Timer struct {
	Name       string `xml:"name,attr"`
	TickPolicy string `xml:"tickpolicy,attr,omitempty"`
	Present    string `xml:"present,attr,omitempty"`
}

type Clock struct {
	Offset string  `xml:"offset,attr"`
	Timers []Timer `xml:"timer"`
}

type Cpu struct {
	Mode string `xml:"mode,attr"`
}

type Features struct {
	Acpi string `xml:"acpi"`
	Apic string `xml:"apic"`
}

type Boot struct {
	Dev string `xml:"dev,attr"`
}

type Type struct {
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
	Value   string `xml:",chardata"`
}

type Os struct {
	Type Type `xml:"type"`
	Boot Boot `xml:"boot"`
}

type Domain struct {
	XMLName       xml.Name `xml:"domain"`
	Type          string   `xml:"type,attr"`
	Name          string   `xml:"name"`
	UUID          string   `xml:"uuid"`
	Memory        int      `xml:"memory"`
	CurrentMemory int      `xml:"currentMemory"`
	VCPU          int      `xml:"vcpu"`
	OS            Os       `xml:"os"`
	Features      Features `xml:"features"`
	CPU           Cpu      `xml:"cpu"`
	Clock         Clock    `xml:"clock"`
	PM            Pm       `xml:"pm"`
	Devices       Devices  `xml:"devices"`
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
    <console type="pty"/>
  </devices>
</domain>`

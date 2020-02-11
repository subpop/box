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

type domainCapabilities struct {
	Path    string `xml:"path"`
	Domain  string `xml:"domain"`
	Machine string `xml:"machine"`
	Arch    string `xml:"arch"`
	VCPU    struct {
		Max string `xml:"max,attr,omitempty"`
	} `xml:"vcpu,omitempty"`
	IOThreads struct {
		Supported string `xml:"supported,attr,omitempty"`
	} `xml:"iothreads,omitempty"`
	OS struct {
		Supported string `xml:"supported,attr,omitempty"`
		Enum      []struct {
			Name  string   `xml:"name,attr,omitempty"`
			Value []string `xml:"value,omitempty"`
		} `xml:"enum,omitempty"`
		Loader struct {
			Supported string `xml:"supported,attr,omitempty"`
			Value     string `xml:"value,omitempty"`
			Enum      []struct {
				Name  string   `xml:"name,attr,omitempty"`
				Value []string `xml:"value,omitempty"`
			} `xml:"enum,omitempty"`
		} `xml:"loader,omitempty"`
	} `xml:"os,omitempty"`
	CPU struct {
		Mode []struct {
			Name      string `xml:"name,attr,omitempty"`
			Supported string `xml:"supported,attr,omitempty"`
			Model     struct {
				Fallback string `xml:"fallback,attr,omitempty"`
				Usable   string `xml:"usable,attr,omitempty"`
				CharData string `xml:",chardata"`
			} `xml:"model,omitempty"`
			Vendor  string `xml:"vendor,omitempty"`
			Feature []struct {
				Policy string `xml:"policy,attr,omitempty"`
				Name   string `xml:"name,attr,omitempty"`
			} `xml:"feature,omitempty"`
		} `xml:"mode,omitempty"`
	} `xml:"cpu,omitempty"`
	Devices struct {
		Disk struct {
			Supported string `xml:"supported,attr,omitempty"`
			Enum      []struct {
				Name  string   `xml:"name,attr,omitempty"`
				Value []string `xml:"value,omitempty"`
			} `xml:"enum,omitempty"`
		} `xml:"disk,omitempty"`
		Graphics struct {
			Supported string `xml:"supported,attr,omitempty"`
			Enum      []struct {
				Name  string   `xml:"name,attr,omitempty"`
				Value []string `xml:"value,omitempty"`
			} `xml:"enum,omitempty"`
		} `xml:"graphics,omitempty"`
		Video struct {
			Supported string `xml:"supported,attr,omitempty"`
			Enum      []struct {
				Name  string   `xml:"name,attr,omitempty"`
				Value []string `xml:"value,omitempty"`
			} `xml:"enum,omitempty"`
		} `xml:"video,omitempty"`
		HostDev struct {
			Supported string `xml:"supported,attr,omitempty"`
			Enum      []struct {
				Name  string   `xml:"name,attr,omitempty"`
				Value []string `xml:"value,omitempty"`
			} `xml:"enum,omitempty"`
		} `xml:"hostdev,omitempty"`
	} `xml:"devices,omitempty"`
	Features struct {
		GIC struct {
			Supported string `xml:"supported,attr,omitempty"`
		} `xml:"gic,omitempty"`
		VMCoreInfo struct {
			Supported string `xml:"supported,attr,omitempty"`
		} `xml:"vmcoreinfo,omitempty"`
		GenID struct {
			Supported string `xml:"supported,attr,omitempty"`
		} `xml:"genid,omitempty"`
		Sev struct {
			Supported string `xml:"supported,attr,omitempty"`
		} `xml:"sev,omitempty"`
	} `xml:"features,omitempty"`
}

type capabilities struct {
	Host struct {
		UUID string `xml:"uuid,omitempty"`
		CPU  struct {
			Arch      string `xml:"arch,omitempty"`
			Model     string `xml:"model,omitempty"`
			Vendor    string `xml:"vendor,omitempty"`
			Microcode struct {
				Version string `xml:"version,attr,omitempty"`
			} `xml:"microcode,omitempty"`
			Counter struct {
				Name      string `xml:"name,attr,omitempty"`
				Frequency string `xml:"frequency,attr,omitempty"`
				Scaling   string `xml:"scaling,attr,omitempty"`
			} `xml:"counter,omitempty"`
			Topology struct {
				Sockets string `xml:"sockets,attr,omitempty"`
				Cores   string `xml:"cores,attr,omitempty"`
				Threads string `xml:"threads,attr,omitempty"`
			} `xml:"topology,omitempty"`
			Features []struct {
				Name string `xml:"name,attr,omitempty"`
			} `xml:"feature,omitempty"`
			Pages []struct {
				Unit string `xml:"unit,attr,omitempty"`
				Size string `xml:"size,attr,omitempty"`
			} `xml:"pages,omitempty"`
		} `xml:"cpu,omitempty"`
		PowerManagement struct {
			SuspendMem    struct{} `xml:"suspend_mem,omitempty"`
			SuspendDisk   struct{} `xml:"suspend_disk,omitempty"`
			SuspendHybrid struct{} `xml:"suspend_hybrid,omitempty"`
		} `xml:"power_management,omitempty"`
		IOMMU struct {
			Support string `xml:"support,attr,omitempty"`
		} `xml:"iommu,omitempty"`
		MigrationFeatures struct {
			Live          struct{} `xml:"live,omitempty"`
			URITransports struct {
				URITransport []string `xml:"uri_transport,omitempty"`
			} `xml:"uri_transports,omitempty"`
		} `xml:"migration_features,omitempty"`
		Topology struct {
			Cells struct {
				Num  string `xml:"num,attr,omitempty"`
				Cell []struct {
					ID     string `xml:"id,attr,omitempty"`
					Memory struct {
						Unit     string `xml:"unit,attr,omitempty"`
						CharData string `xml:",chardata"`
					} `xml:"memory,omitempty"`
					Pages []struct {
						Unit     string `xml:"unit,attr,omitempty"`
						Size     string `xml:"size,attr,omitempty"`
						CharData string `xml:",chardata"`
					} `xml:"pages,omitempty"`
					Distances struct {
						Sibling struct {
							ID    string `xml:"id,attr,omitempty"`
							Value string `xml:"value,attr,omitempty"`
						} `xml:"sibling,omitempty"`
					} `xml:"distances,omitempty"`
					CPUs struct {
						Num string `xml:"num,attr,omitempty"`
						CPU []struct {
							ID       string `xml:"id,attr,omitempty"`
							SocketID string `xml:"socket_id,attr,omitempty"`
							CoreID   string `xml:"core_id,attr,omitempty"`
							Siblings string `xml:"siblings,attr,omitempty"`
						} `xml:"cpu,omitempty"`
					} `xml:"cpus,omitempty"`
				} `xml:"cell,omitempty"`
			} `xml:"cells,omitempty"`
		} `xml:"topology,omitempty"`
		Cache struct {
			Bank struct {
				ID    string `xml:"id,attr,omitempty"`
				Level string `xml:"level,attr,omitempty"`
				Type  string `xml:"type,attr,omitempty"`
				Size  string `xml:"size,attr,omitempty"`
				Unit  string `xml:"unit,attr,omitempty"`
				CPUS  string `xml:"cpus,attr,omitempty"`
			} `xml:"bank,omitempty"`
		} `xml:"cache,omitempty"`
		SecModel struct {
			Model     string `xml:"model,omitempty"`
			DOI       string `xml:"doi,omitempty"`
			BaseLabel []struct {
				Type     string `xml:"type,attr,omitempty"`
				CharData string `xml:",chardata"`
			} `xml:"baselabel,omitempty"`
		} `xml:"secmodel,omitempty"`
	} `xml:"host,omitempty"`
	Guest []struct {
		OSType string `xml:"os_type,omitempty"`
		Arch   struct {
			Name     string `xml:"name,attr,omitempty"`
			WordSize string `xml:"wordsize,omitempty"`
			Emulator string `xml:"emulator,omitempty"`
			Machine  []struct {
				MaxCPUs   string `xml:"maxCpus,attr,omitempty"`
				Canonical string `xml:"canonical,attr,omitempty"`
				CharData  string `xml:",chardata"`
			} `xml:"machine,omitempty"`
			Domain []struct {
				Type string `xml:"type,attr,omitempty"`
			} `xml:"domain,omitempty"`
		} `xml:"arch,omitempty"`
		Features struct {
			CPUSelection struct{} `xml:"cpuselection,omitempty"`
			DeviceBoot   struct{} `xml:"deviceboot,omitempty"`
			DiskSnapshot struct{} `xml:"disksnapshot,omitempty"`
			ACPI         struct {
				Default string `xml:"default,attr,omitempty"`
				Toggle  string `xml:"toggle,attr,omitempty"`
			} `xml:"acpi,omitempty"`
			APIC struct {
				Default string `xml:"default,attr,omitempty"`
				Toggle  string `xml:"toggle,attr,omitempty"`
			} `xml:"apic,omitempty"`
			PAE    struct{} `xml:"pae,omitempty"`
			NonPAE struct{} `xml:"nonpae,omitempty"`
		} `xml:"features,omitempty"`
	} `xml:"guest,omitempty"`
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

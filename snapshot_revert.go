package vm

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
)

// SnapshotRevert discards the current domain state, reverting it to snapshotName.
func SnapshotRevert(domainName, snapshotName string) error {
	var err error

	conn, err := libvirt.NewConnect("")
	if err != nil {
		return err
	}

	dom, err := conn.LookupDomainByName(domainName)
	if err != nil {
		return err
	}
	defer dom.Free()

	snapshot, err := dom.SnapshotLookupByName(snapshotName, 0)
	if err != nil {
		return err
	}

	err = snapshot.RevertToSnapshot(0)
	if err != nil {
		return err
	}

	fmt.Println("Domain reverted to snapshot " + snapshotName)

	return nil
}

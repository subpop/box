package vm

import (
	"fmt"
	"os"

	"github.com/libvirt/libvirt-go"
)

// A MissingPositionalArgErr occurs when a command is invoked without a required
// positional argument.
type MissingPositionalArgErr struct {
	name string
}

func (e MissingPositionalArgErr) Error() string {
	return fmt.Sprintf("error: %v is required", e.name)
}

// ErrDomainNameRequired represents a missing domain name argument.
var ErrDomainNameRequired = MissingPositionalArgErr{
	name: "domain name",
}

// ErrImageNameRequired represents a missing image name argument.
var ErrImageNameRequired = MissingPositionalArgErr{
	name: "image name",
}

// ErrTemplateNameRequired represents a missing template name argument.
var ErrTemplateNameRequired = MissingPositionalArgErr{
	name: "template name",
}

// ErrURLOrPathRequired represents a missing URL or path argument.
var ErrURLOrPathRequired = MissingPositionalArgErr{
	name: "URL or path",
}

// LogErrorAndExit logs err and exits with a non-zero exit code.
func LogErrorAndExit(err error) {
	switch err.(type) {
	case libvirt.Error:
		e := err.(libvirt.Error)
		fmt.Fprintf(os.Stderr, "%v\n", e.Message)
		os.Exit(int(e.Code))
	default:
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

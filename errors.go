package vm

import (
	"fmt"
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

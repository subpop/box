package vm

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uiprogress"
)

type byteCounter struct {
	bar    *uiprogress.Bar
	total  uint64
	prefix string
}

func (b *byteCounter) Write(p []byte) (int, error) {
	n := len(p)
	b.total += uint64(n)
	if b.bar != nil {
		b.bar.Set(b.bar.Current() + n)
	} else {
		fmt.Printf("\r%s", strings.Repeat(" ", 35))
		fmt.Printf("\r%s%s", b.prefix, humanize.Bytes(b.total))
	}

	return n, nil
}

func (b *byteCounter) Close() error {
	fmt.Println()

	return nil
}

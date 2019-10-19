package vm

// byteCounter is a Writer that invokes a callback each time Write is called.
// It does not write the byte slice passed to it anywhere, so it must be used
// with an io.TeeReader in order to actually write data anywhere meaningful.
type byteCounter struct {
	onWrite func(b []byte)
	onClose func()
}

func (b *byteCounter) Write(buf []byte) (int, error) {
	b.onWrite(buf)
	return len(buf), nil
}

func (b *byteCounter) Close() error {
	b.onClose()
	return nil
}

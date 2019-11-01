package vm

// ImageGet downloads rawurl and prepares it for use as a backing disk image. If
// quiet is true, no progress is printed to stdout.
func ImageGet(rawurl string, quiet bool) error {
	var err error

	filePath, err := download(rawurl, quiet)
	if err != nil {
		return err
	}

	_, err = inspect(filePath, quiet)
	if err != nil {
		return err
	}

	return nil
}

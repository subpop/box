package vm

// ImageGet downloads rawurl and prepares it for use as a backing disk image.
func ImageGet(rawurl string) error {
	var err error

	filePath, err := download(rawurl)
	if err != nil {
		return err
	}

	err = inspect(filePath)
	if err != nil {
		return err
	}

	return nil
}

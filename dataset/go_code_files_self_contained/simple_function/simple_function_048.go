package main

func prepareExecutable() error {
	path, err := downloadFromUrl(GetDownloadURL(), base+"/$V", thrustVersion)
	if err != nil {
		return err
	}

	return UnzipExecutable(path)
}

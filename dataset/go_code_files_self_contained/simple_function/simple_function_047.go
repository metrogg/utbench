package main

func getOSImageVersion() (string, error) {
	productName, err := getCurrentVersionVal("ProductName")
	if err != nil {
		return "", nil
	}

	return productName, nil
}

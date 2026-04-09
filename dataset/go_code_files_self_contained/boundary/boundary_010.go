package main

func Width() (uint, error) {
	output, err := size()
	if err != nil {
		return 0, err
	}
	_, width, err := parse(output)
	return width, err
}

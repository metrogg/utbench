package main

func Height() (uint, error) {
	output, err := size()
	if err != nil {
		return 0, err
	}
	height, _, err := parse(output)
	return height, err
}

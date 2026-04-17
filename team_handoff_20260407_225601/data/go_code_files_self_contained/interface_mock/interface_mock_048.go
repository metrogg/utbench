package main

func GzipAssetNames() []string {
	names := make([]string, 0, len(_gzipbindata))
	for name := range _gzipbindata {
		names = append(names, name)
	}
	return names
}

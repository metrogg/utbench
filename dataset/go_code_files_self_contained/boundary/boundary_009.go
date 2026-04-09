package main

func genArchive() ([]string, error) {
	switch *flagArchiveType {
	case "linux", "darwin", "windows":
		genBinaries(*flagArchiveType)
		return []string{packBinaries(*flagArchiveType)}, nil
	case "src":
		return []string{zipSource()}, nil
	default:
		return genAll()
	}
}

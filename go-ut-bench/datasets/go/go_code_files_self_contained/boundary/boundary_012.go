package main

func GetPluginInstallPrefixes() ([]string, error) {
	primaryPluginInstallDir, err := GetPrimaryPluginsInstallDir()
	if err != nil {
		return nil, err
	}
	return []string{primaryPluginInstallDir}, nil
}

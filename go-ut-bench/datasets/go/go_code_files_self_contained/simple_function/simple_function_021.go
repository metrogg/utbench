package main

func SelinuxEnabled() bool {
	if selinuxEnabledChecked {
		return selinuxEnabled
	}
	selinuxEnabledChecked = true
	if fs := getSelinuxMountPoint(); fs != "" {
		if con, _ := Getcon(); con != "kernel" {
			selinuxEnabled = true
		}
	}
	return selinuxEnabled
}

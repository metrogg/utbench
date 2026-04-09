package main

import (
	"io/ioutil"
	"path/filepath"
)

func EnableIPv6Forwarding() error {
	ip6ConfPath := "/proc/sys/net/ipv6/conf/"
	device := "all"
	forwarding := "forwarding"
	forwardingOn := "1"
	path := filepath.Join(ip6ConfPath, device, forwarding)
	return ioutil.WriteFile(path, []byte(forwardingOn), 0644)
}

package main

import (
	"encoding/json"
	"strings"
)

func InstanceTags() ([]string, error) {
	var s []string
	j, err := Get("instance/tags")
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(strings.NewReader(j)).Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

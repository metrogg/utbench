package main

import (
	"fmt"
)

func Who() (int, int) {
	// subs := make(map[string]bool) // no duplicates
	subs := []string{}

	// get tags all clients subscribed to
	for i := range subscribers {
		subs = append(subs, fmt.Sprint(subscribers[i].id))
	}

	return len(subs), int(uid)
}

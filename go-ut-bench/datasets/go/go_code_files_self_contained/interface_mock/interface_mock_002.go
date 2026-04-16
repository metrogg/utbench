package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 2 && strings.ToLower(os.Args[1]) == "server" {
		ServerLoop()
	} else if len(os.Args) == 2 && strings.ToLower(os.Args[1]) == "client" {
		ClientLoop()
	} else {
		fmt.Println("wsserver [server|client]")
	}
}

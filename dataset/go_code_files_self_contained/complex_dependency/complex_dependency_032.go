package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/_ah/background", handleBackground)

	sc := make(chan send)
	rc := make(chan recv)
	sendc, recvc = sc, rc
	go matchmaker(sc, rc)
}

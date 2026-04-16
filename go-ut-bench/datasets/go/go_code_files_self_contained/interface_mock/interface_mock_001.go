package main

import (
	"io"
	"net/http"
)

func HTTPHandlerGetLevel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, GetLevel().String())
	})
}

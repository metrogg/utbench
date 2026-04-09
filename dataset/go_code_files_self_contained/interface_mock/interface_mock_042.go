package main

import (
	"net/http"
)

func New() *Context {
	req := createRequest()
	res := createResponse(req)
	cli := &http.Client{Transport: http.DefaultTransport}
	return &Context{Request: req, Response: res, Client: cli}
}

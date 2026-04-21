package main

import (
	"net/http"
	"net/url"
	"os"
)

func transportFromEnv() http.RoundTripper {
	for _, v := range []string{"HTTPS_PROXY", "HTTP_PROXY"} {
		if proxy := os.Getenv(v); proxy != "" {
			if u, err := url.Parse(proxy); err == nil {
				return &http.Transport{
					Proxy: http.ProxyURL(u),
				}
			}
		}
	}
	return nil
}

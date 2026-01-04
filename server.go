package main

import "net/http"

func newServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Handler: router,
		Addr:    addr,
	}
}

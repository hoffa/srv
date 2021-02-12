package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func listenAndServe(handler http.Handler, addr string, timeout time.Duration) error {
	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return srv.ListenAndServe()
}

func listenAndServeTLS(handler http.Handler, addr string, timeout time.Duration, host, certDir string) error {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
	}
	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		TLSConfig:    m.TLSConfig(),
	}
	return srv.ListenAndServeTLS("", "")
}

func main() {
	addr := flag.String("addr", ":443", "address the server listens to")
	dir := flag.String("dir", ".", "directory to serve")
	tlsHost := flag.String("https", "", "domain name of the server (TLS is disabled if empty)")
	certDir := flag.String("cert", ".", "directory for storing TLS certificates")
	timeout := flag.Duration("timeout", time.Minute, "server timeout")
	flag.Parse()

	handler := http.FileServer(http.Dir(*dir))
	if *tlsHost == "" {
		fmt.Printf("Serving %v on %v\n", *dir, *addr)
		panic(listenAndServe(handler, *addr, *timeout))
	} else {
		fmt.Printf("Serving %v on %v with TLS (%v)\n", *dir, *addr, *tlsHost)
		panic(listenAndServeTLS(handler, *addr, *timeout, *tlsHost, *certDir))
	}
}

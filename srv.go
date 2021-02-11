package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"time"
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
	addr := flag.String("addr", ":443", "")
	dir := flag.String("dir", ".", "")
	tlsHost := flag.String("https", "", "")
	certDir := flag.String("cert", ".", "")
	timeout := flag.Duration("timeout", time.Minute, "")
	flag.Parse()

	handler := http.FileServer(http.Dir(*dir))
	if *tlsHost == "" {
		fmt.Printf("addr=%v dir=%v timeout=%v", *addr, *dir, *timeout)
		panic(listenAndServe(handler, *addr, *timeout))
	} else {
		fmt.Printf("addr=%v dir=%v timeout=%v https=%v cert=%v", *addr, *dir, *timeout, *tlsHost, *certDir)
		panic(listenAndServeTLS(handler, *addr, *timeout, *tlsHost, *certDir))
	}
}

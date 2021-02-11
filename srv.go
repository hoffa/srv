package main

import (
	"flag"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"time"
)

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
	log.Fatal(listenAndServeTLS(handler, *addr, *timeout, *tlsHost, *certDir))
}
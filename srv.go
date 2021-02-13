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

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	dir := flag.String("dir", ".", "directory to serve")
	certDir := flag.String("cert", ".", "directory for storing TLS certificates")
	timeout := flag.Duration("timeout", time.Minute, "server timeout")
	flag.Parse()
	host := flag.Arg(0)

	handler := http.FileServer(http.Dir(*dir))
	fmt.Printf("Serving %v on %v\n", *dir, host)
	go func() {
		panic(listenAndServe(http.HandlerFunc(redirectTLS), ":80", *timeout))
	}()
	panic(listenAndServeTLS(handler, ":443", *timeout, host, *certDir))
}

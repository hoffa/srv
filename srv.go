package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func main() {
	var domain string
	var addr string
	var dir string
	var certDir string
	var timeout time.Duration

	flag.StringVar(&domain, "n", "", "domain name (required)")
	flag.StringVar(&dir, "d", "", "directory to serve (required)")
	flag.StringVar(&certDir, "c", "", "directory for storing TLS certificates (required)")
	flag.StringVar(&addr, "p", ":443", "HTTPS address")
	flag.DurationVar(&timeout, "t", 10*time.Second, "timeout")
	flag.Parse()

	if domain == "" || dir == "" || certDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}
	srv := &http.Server{
		Handler:      loggingHandler(http.FileServer(http.Dir(dir))),
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		TLSConfig:    m.TLSConfig(),
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

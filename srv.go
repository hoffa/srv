package main

import (
	"flag"
	"io/ioutil"
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

func listenAndServeTLS(addr string, timeout time.Duration, host, certDir string, handler http.Handler) error {
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

func tempDir() string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return dir
}

func main() {
	var arg struct {
		domain  string
		addr    string
		dir     string
		certDir string
		timeout time.Duration
	}

	flag.StringVar(&arg.domain, "n", "", "domain name (required)")
	flag.StringVar(&arg.addr, "p", ":443", "HTTPS address")
	flag.StringVar(&arg.dir, "d", ".", "directory to serve")
	flag.StringVar(&arg.certDir, "c", "", "directory to store TLS certificates (temporary directory if not set)")
	flag.DurationVar(&arg.timeout, "t", 10*time.Second, "timeout")
	flag.Parse()

	handler := loggingHandler(http.FileServer(http.Dir(arg.dir)))

	if arg.domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	if arg.certDir == "" {
		arg.certDir = tempDir()
	}
	log.Fatal(listenAndServeTLS(arg.addr, arg.timeout, arg.domain, arg.certDir, handler))
}

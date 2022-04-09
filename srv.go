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

func listenAndServe(addr string, timeout time.Duration, handler http.Handler) error {
	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return srv.ListenAndServe()
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
		host      string
		httpsAddr string
		dir       string
		certDir   string
		timeout   time.Duration
	}

	flag.StringVar(&arg.host, "n", "", "domain name (required)")
	flag.StringVar(&arg.httpsAddr, "p", ":443", "HTTPS address")
	flag.StringVar(&arg.dir, "d", ".", "directory to serve")
	flag.StringVar(&arg.certDir, "c", "", "directory to store TLS certificates (temporary directory if not set)")
	flag.DurationVar(&arg.timeout, "t", time.Minute, "timeout")
	flag.Parse()

	handler := loggingHandler(http.FileServer(http.Dir(arg.dir)))

	if arg.host == "" {
		flag.Usage()
		os.Exit(1)
	}

	if arg.certDir == "" {
		arg.certDir = tempDir()
	}
	log.Fatal(listenAndServeTLS(arg.httpsAddr, arg.timeout, arg.host, arg.certDir, handler))
}

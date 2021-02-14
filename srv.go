package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

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
		httpsPort int
		httpPort  int
		dir       string
		certDir   string
		timeout   time.Duration
	}

	flag.StringVar(&arg.host, "n", "", "domain name (required)")
	flag.IntVar(&arg.httpsPort, "p", 443, "HTTPS port")
	flag.IntVar(&arg.httpPort, "q", 80, "HTTP redirect port")
	flag.StringVar(&arg.dir, "d", ".", "directory to serve")
	flag.StringVar(&arg.certDir, "c", "", "directory to store TLS certificates (temporary directory if empty)")
	flag.DurationVar(&arg.timeout, "t", time.Minute, "timeout")
	flag.Parse()

	if arg.host == "" {
		flag.Usage()
		os.Exit(1)
	}

	if arg.certDir == "" {
		arg.certDir = tempDir()
	}

	fmt.Printf("%+v\n", arg)

	httpsPort := strconv.Itoa(arg.httpsPort)
	httpPort := strconv.Itoa(arg.httpPort)

	go func() {
		panic(listenAndServe(":"+httpPort, arg.timeout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := "https://" + arg.host + ":" + httpsPort + r.RequestURI
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		})))
	}()

	handler := http.FileServer(http.Dir(arg.dir))
	panic(listenAndServeTLS(":"+httpsPort, arg.timeout, arg.host, arg.certDir, handler))
}

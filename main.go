package main

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const LOG_LEVEL = "LOG_LEVEL"

func main() {
	if os.Getenv(LOG_LEVEL) != "" {
		logLevel, err := log.ParseLevel(os.Getenv(LOG_LEVEL))
		if err != nil {
			log.Warnf("Unable to parse LOG_LEVEL value using default: %v", err)
		} else {
			log.SetLevel(logLevel)
		}
	}
	listenAddr := flag.String("listen", "0.0.0.0:80", "host:port to listen for oncoming requests")
	forwardAddr := flag.String("forward", "127.0.0.1:3000", "host:port to forward modified requests to")
	headers := make(headerFlags)
	flag.Var(&headers, "header", "Headers to inject in name=value format")
	flag.Parse()

	checkAddr(*listenAddr, "listen")
	checkAddr(*forwardAddr, "forward")
	forwardBaseUrl := &url.URL{
		Scheme: "http",
		Host:   *forwardAddr,
		Path:   "/",
	}

	log.Infof("Listening on %s", *listenAddr)
	log.Infof("Forwarding to %s", forwardBaseUrl)
	log.Infof("Injecting headers %v", headers)

	rp := httputil.NewSingleHostReverseProxy(forwardBaseUrl)

	s := &http.Server{
		Addr: *listenAddr,
		Handler: &headerHandler{
			forwardAddr:     *forwardAddr,
			injectedHeaders: headers,
			reverseProxy:    rp,
		},
	}

	log.Fatal(s.ListenAndServe())
}

type headerHandler struct {
	forwardAddr     string
	injectedHeaders headerFlags
	reverseProxy    *httputil.ReverseProxy
}

func (th *headerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Processing request %s", r.RequestURI)
	for key, value := range th.injectedHeaders {
		r.Header.Add(key, value)
	}
	th.reverseProxy.ServeHTTP(w, r)
}

type headerFlags map[string]string

func (i *headerFlags) String() string {
	return "my string representation"
}

func (i *headerFlags) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return errors.New(fmt.Sprintf("headers must be in the name=value format"))
	}
	(*i)[parts[0]] = parts[1]
	return nil
}

func checkAddr(addr, flagName string) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "invalid value %q for flag -%s: expect host:port format\n", addr, flagName)
		flag.Usage()
		os.Exit(1)
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:80", "host:port to listen for oncoming requests")
	forwardAddr := flag.String("forward", "127.0.0.1:3000", "host:port to forward modified requests to")
	headers := make(headerFlags)
	flag.Var(&headers, "header", "Headers to inject in name=value format")
	flag.Parse()

	checkAddr(*listenAddr, "listen")
	checkAddr(*forwardAddr, "forward")

	log.Printf("Listening on %s", *listenAddr)
	log.Printf("Forwarding to %s", *forwardAddr)
	log.Printf("Injecting headers %v", headers)

	s := &http.Server{
		Addr: *listenAddr,
		Handler: &headerHandler{
			forwardAddr:     *forwardAddr,
			injectedHeaders: headers,
		},
	}

	log.Fatal(s.ListenAndServe())
}

type headerHandler struct {
	forwardAddr     string
	injectedHeaders headerFlags
}

func (th *headerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	r.Host = th.forwardAddr
	url := mapUrl(r)
	for key, value := range th.injectedHeaders {
		r.Header.Add(key, value)
	}
	nr := cloneRequest(r, url)
	resp, err := client.Do(nr)
	if err != nil {
		errorResponse(w, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Write([]byte(body))
}

func mapUrl(r *http.Request) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   r.RequestURI,
	}
}

func cloneRequest(r *http.Request, url *url.URL) *http.Request {
	return &http.Request{
		Method:     r.Method,
		URL:        url,
		Header:     r.Header.Clone(),
		Body:       r.Body,
		Host:       "",
		RemoteAddr: r.RemoteAddr,
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintf("Internal Error: %v", err)))
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

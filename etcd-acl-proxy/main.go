package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/proxy"
)

var (
	addr  = flag.String("addr", ":7777", "http port")
)

func ro(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && strings.HasPrefix(req.URL.Path, "/v2/keys/read-only") {
			log.Printf("denied a %s on read-only", req.Method)
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		next.ServeHTTP(w, req)
	}
}

func wo(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		println(req.URL.Path)
		if req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/v2/keys/write-only") {
			log.Print("denied a GET on write-only")
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		next.ServeHTTP(w, req)
	}
}

func main() {
	pt := &http.Transport{
		// timeouts taken from http.DefaultTransport
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	ma, _ := client.NewMembersAPI(pt, "http://localhost:7001", 15 * time.Second)
	// TODO(philips): persist to disk
	memURLs := []string{"http://localhost:4001"}
	uf := func() []string {
		mems, err := ma.List()
		if err != nil {
			log.Print("unable to list members.")
			return memURLs
		}
		for _, k := range mems {
			memURLs = append(memURLs, k.ClientURLs...)
		}
		return memURLs
	}

	ph := proxy.NewHandler(pt, uf)
	ph = http.HandlerFunc(wo(ph))
	ph = http.HandlerFunc(ro(ph))

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("etcd: proxy listening for client requests on ", *addr)
	log.Fatal(http.Serve(l, ph))
}

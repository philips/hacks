package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/philips/hacks/host-info/Godeps/_workspace/src/github.com/dgryski/go-identicon"
)

const (
	home = `<html>
	<head>
		<title>{{.Hostname}} Info</title>
		<link rel="shortcut icon" href="/icons/{{.Hostname}}.png" type="image/x-icon" />
	</head>
	<body>
		<img style="float: right; padding: 30px; position: relative; width: 30%; display: inline;" src="/icons/{{.Hostname}}.png" />
		<h1><p>You are visitor: {{.VisitCount}}</p><h1>
		<h1>{{.Hostname}}</h1>
		<h1>{{.Zone}}</h1>
		<h2>Networking</h2>
		{{range .Interfaces}}
			<h3>{{.Name}}</h3>
			{{range .Addrs}}
				<ul>
					<li>{{.String}}</li>
				</ul>
			{{end}}
		{{end}}
	</body>
</html>`
)

type Host struct {
	VisitCount uint64
	Hostname   string
	Zone       string
	Interfaces []net.Interface
}

func NewHost() Host {
	h := Host{}

	hostname, _ := os.Hostname()

	h.Hostname = hostname
	interfaces, _ := net.Interfaces()
	h.Interfaces = interfaces
	h.Zone = getZone()
	return h
}

func (h *Host) Write(w io.Writer) {
	atomic.AddUint64(&h.VisitCount, 1)
	tmpl, err := template.New("index").Parse(home)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, h)
	if err != nil {
		panic(err)
	}
}

func iconHandler(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/icons/")
	args = args[1:]

	if len(args) != 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := args[0]

	// support jpg too?
	if !strings.HasSuffix(item, ".png") {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item = strings.TrimSuffix(item, ".png")

	key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	icon := identicon.New7x7(key)

	log.Printf("creating identicon for '%s'\n", item)

	data := []byte(item)
	pngdata := icon.Render(data)

	w.Header().Set("Content-Type", "image/png")
	w.Write(pngdata)

	return
}

const zoneUrl string = "http://metadata.google.internal/computeMetadata/v1/instance/zone"

func getZone() string {
	timeout := time.Duration(10 * time.Millisecond)
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", zoneUrl, nil)
	req.Header.Add("Metadata-Flavor", "Google")
	response, err := client.Do(req)
	if err != nil {
		return ("Unknown Zone")
	} else {
		defer response.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		zoneSplit := strings.Split(buf.String(), "/")
		return zoneSplit[len(zoneSplit)-1]
	}
}

func main() {
	addr := ":8080"

	hostInfo := NewHost()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request from %v\n", r.RemoteAddr)
		hostInfo.Write(w)
	})
	http.HandleFunc("/icons/", iconHandler)
	log.Printf("listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

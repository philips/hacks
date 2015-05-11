package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	testList   []Events
	callEvents []Events
)

func init() {
	testList = []Events{
		Events{
			[]Event{
				Event{unixMS(time.Now()), unixMS(time.Now().Add(time.Minute)), "a"},
				Event{unixMS(time.Now().Add(time.Minute * 2)), unixMS(time.Now().Add(time.Minute * 3)), "b"},
			},
			"foo",
		},
		Events{
			[]Event{
				Event{unixMS(time.Now().Add(time.Minute * 4)), unixMS(time.Now().Add(time.Minute * 5)), "c"},
			},
			"baz",
		},
		Events{
			[]Event{
				Event{unixMS(time.Now().Add(time.Minute * 8)), unixMS(time.Now().Add(time.Minute * 12)), "d"},
			},
			"bar",
		},
	}

}

func main() {
	parseData()
	http.Handle("/", http.FileServer(http.Dir("root")))
	http.HandleFunc("/api", myhandler)

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
}

type Events struct {
	Times []Event `json:"times"`
	Label string  `json:"label"`
}

type Event struct {
	StartingTime int64  `json:"starting_time"`
	EndingTime   int64  `json:"ending_time"`
	Label        string `json:"label"`
}

func unixMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func parseData() {
	raw, err := os.Open("rawdata.txt")
	if err != nil {
		log.Printf("%v", err)
	}

	data := []url.Values{}
	scanner := bufio.NewScanner(raw)
	for scanner.Scan() {
		q, err := url.ParseQuery(scanner.Text())
		if err != nil {
			log.Printf("%v", err)
		}

		data = append(data, q)
	}

	callEvents = []Events{}

	for _, d := range data {
		csid := d.Get("CallSid")[:4]
		status := d.Get("CallStatus")[:4]
		t, err := time.Parse(time.RFC1123Z, d.Get("Timestamp"))
		println(csid, status, t.String())

		if err != nil {
			log.Printf("%v", err)
		}

		event := Event{unixMS(t), unixMS(t.Add(time.Second)), status}

		found := false
		for i, exist := range callEvents {
			if exist.Label != csid {
				continue
			}
			callEvents[i].Times = append(exist.Times, event)
			found = true
		}
		if found == false {
			e := Events{
				[]Event{event},
				csid,
			}

			callEvents = append(callEvents, e)
		}

	}

}

func myhandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("got api request")
	b, err := json.Marshal(callEvents)
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Fprintf(w, string(b))
}

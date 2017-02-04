package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)


func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestAuth(t *testing.T) {
	tests := []struct {
		Method string
		Path string
		StatusCode int
	}{
		{"PUT", "/accessible/baz", 200},
		{"PUT", "/inaccessible",  403},
		{"GET", "/accessible", 200},
		{"GET", "/inaccessible", 403},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	client := &http.Client{}

	for i, test := range tests {
		req := newRequest(test.Method, ts.URL+test.Path)
		resp, err := client.Do(req)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != test.StatusCode {
			t.Fatalf("%d: wrong code, got %d want %d", i, resp.StatusCode, test.StatusCode)
		}
	}
}

package main

import (
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name   string
		wfound bool
	}{
		{"invalid-always.coreos.com.", false},
		{"www.coreos.com.", true},
	}

	for i, tt := range tests {
		msg, _, found := find(tt.name)
		if tt.wfound != found {
			t.Errorf("#%d: name = %s, msg = %s, want %b", i, tt.name, msg, tt.wfound)
		}
	}
}

package main

import "testing"

func TestQueryDNS(t *testing.T) {

	ip := QueryDNS("tomasen.org", "")

	if string(ip) != "72.14.188.94" {
		t.Fatal()
	}
}

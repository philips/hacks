package main

import (
	"fmt"
	"os"
	"time"

	"github.com/miekg/dns"
)

func rrHeaderFmt(hdr dns.RR_Header) string {
	return fmt.Sprintf("'%s' %s ttl: %d", hdr.Name, dns.Type(hdr.Rrtype).String(), hdr.Ttl)
}

func aFmt(record *dns.A) string {
	return fmt.Sprintf("'%s' ttl: %d", record.A.String(), record.Hdr.Ttl)
}

func find(host string) (status string, ttl uint32, found bool) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{host, dns.TypeA, dns.ClassINET}

	c := new(dns.Client)
	in, _, err := c.Exchange(m1, "8.8.8.8:53")

	if err != nil {
		return fmt.Sprintf("Error on query for '%s': %v", host, err), 0, false
	}

	if len(in.Answer) == 0 {
		var longest dns.RR_Header
		for _, n := range in.Ns {
			if n.Header().Ttl > longest.Ttl {
				longest = *n.Header()
			}
		}
		return fmt.Sprintf("No answer on query for '%s'. Waiting for timeout of authoritative record %s", host, rrHeaderFmt(longest)), longest.Ttl, false
	}

	return fmt.Sprintf("Got answer on query for '%s': %v", host, aFmt(in.Answer[0].(*dns.A))), 0, true
}

func main() {
	for {
		msg, ttl, found := find("tec6.a.ifup.org.")
		fmt.Printf(msg + "\n")
		if found {
			os.Exit(0)
		}
		sleep := time.Second
		if ttl != 0 {
			sleep = time.Duration(ttl/2) * time.Second
		}
		time.Sleep(sleep)
	}
}

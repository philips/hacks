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

func newMsg(host string, qClass uint16) *dns.Msg {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{host, qClass, dns.ClassINET}
	return m1
}

func makeRequest(ns string, msg *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	resp, _, err := c.Exchange(msg, ns)

	if err != nil {
		return nil, fmt.Errorf("Error on query for '%s': %v", ns, err)
	}
	return resp, nil
}

func find(host string) (status string, ttl uint32, found bool) {
	// get SOA record
	reqNS := newMsg(host, dns.TypeSOA)
	resp, err := makeRequest("8.8.8.8:53", reqNS)
	if err != nil {
		return fmt.Sprintf("could not retrieve SOA record: %v", err), 0, false
	} else if len(resp.Ns) < 1 {
		return "no SOA record found", 0, false
	}

	ns := resp.Ns[0].(*dns.SOA).Ns

	// resolve authoritative nameserver
	nsReqA := newMsg(ns, dns.TypeA)
	resp, err = makeRequest("8.8.8.8:53", nsReqA)
	if err != nil {
		return fmt.Sprintf("could not resolve authoritative ns: %v", err), 0, false
	} else if len(resp.Answer) < 1 {
		return "Authoratiative NS domain does not have an A record", 0, false
	}

	// request A record from authoritative name server
	reqA := newMsg(host, dns.TypeA)
	resp, err = makeRequest(resp.Answer[0].(*dns.A).A.String()+":53", reqA)
	if err != nil {
		return err.Error(), 0, false
	} else if len(resp.Answer) < 1 {
		return "Domain does not have an A record", 0, false
	}

	return fmt.Sprintf("Got answer on query for '%s': %v", host, aFmt(resp.Answer[0].(*dns.A))), 0, true
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

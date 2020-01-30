package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"strconv"
)

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		log.Println("Serving DNS request for " + domain)
		if domain == "local.nacdlow.com." {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("127.0.0.1"),
			})

		} else {
			msg.Answer = getAnswer(domain)
		}
	}
	w.WriteMsg(&msg)
}

var (
	c *dns.Client = new(dns.Client)
)

func getAnswer(dom string) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dom), dns.TypeA)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53"))
	if r == nil {
		log.Fatalf("Error while getting external request: %s\n", err.Error())
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Fatalln("Invalid answer while getting external request!")
	}
	return r.Answer
}

func main() {
	log.Println("Starting DNS server on port 53!")
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}

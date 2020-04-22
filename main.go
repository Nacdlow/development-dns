package main

import (
	"fmt"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
)

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered handler", r)
		}
	}()

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
				A:   net.ParseIP(os.Args[1]),
			})

		} else {
			msg.Answer = getAnswer(domain)
		}
	}
	w.WriteMsg(&msg)
}

var (
	c      *dns.Client = new(dns.Client)
	domain             = "local.nacdlow.com."
	ip     string
	//extDNS = "137.195.151.105"
	extDNS = "8.8.8.8"
)

func getAnswer(dom string) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dom), dns.TypeA)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort(extDNS, "53"))
	if r == nil {
		log.Printf("Error while getting external request: %s\n", err.Error())
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Println("Invalid answer while getting external request!")
	}
	return r.Answer
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass your computer's IP as an argument!")
		os.Exit(-1)
	}
	ip = os.Args[1]
	log.Infof("This server will resolve '%s' to '%s'!\n", domain, ip)
	log.Infof("Other queries will be forwarded to %s\n", extDNS)
	log.Infoln("Starting DNS server on port 53!")
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\nAre you running as root?", err.Error())
	}
}

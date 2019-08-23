package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

func handlerQuery(m *dns.Msg, w dns.ResponseWriter) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			remoteAddr := w.RemoteAddr().(*net.UDPAddr).IP
			remotePort := w.RemoteAddr().(*net.UDPAddr).Port
			if *ttl == -1 {
				*ttl = remotePort
			}
			log.Printf("A: %s, TTL:%d\n", remoteAddr, *ttl)
			rr, err := dns.NewRR(fmt.Sprintf("%s %d IN A %s", q.Name, *ttl, remoteAddr))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeAAAA:
			log.Printf("Query for %s\n", q.Name)
			remoteAddr := w.RemoteAddr().(*net.UDPAddr).IP
			remotePort := w.RemoteAddr().(*net.UDPAddr).Port
			if *ttl == -1 {
				*ttl = remotePort
			}
			log.Printf("A: %s, TTL:%d\n", remoteAddr, *ttl)
			rr, err := dns.NewRR(fmt.Sprintf("%s %d IN AAAA %s", q.Name, *ttl, remoteAddr))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func gelwormHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		handlerQuery(m, w)
	}
	w.WriteMsg(m)
}

var (
	port = flag.Int("port", 15353, "Run DNS port")
	zone = flag.String("zone", ".", "Run DNS zone")
	host = flag.String("host", "0.0.0.0", "Run DNS host")
	ttl = flag.Int("ttl", -1, "DNS TTL")
)

func main() {
	flag.Parse()

	dns.HandleFunc(*zone, gelwormHandler)
	server := &dns.Server{Addr: *host + ":" + strconv.Itoa(*port), Net: "udp"}
	log.Printf("%s zone DNS Server run %s:%d\n", *zone, *host, *port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start: %s\n", err.Error())
	}
}

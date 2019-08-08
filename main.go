package main

import (
	"strconv"
	"log"
	"fmt"
	"net"

	"github.com/miekg/dns"
)


func handlerQuery(m *dns.Msg, w dns.ResponseWriter){
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			remoteAddr := w.RemoteAddr().(*net.UDPAddr).IP
			log.Println(remoteAddr)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, remoteAddr))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		case dns.TypeAAAA:
			log.Printf("Query for %s\n", q.Name)
			remoteAddr := w.RemoteAddr().(*net.UDPAddr).IP
			log.Println(remoteAddr)
			rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, remoteAddr))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func gelwormHandler(w dns.ResponseWriter, r *dns.Msg){
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	log.Println(w.LocalAddr())
	log.Println(w.RemoteAddr())

	switch r.Opcode {
	case dns.OpcodeQuery:
		handlerQuery(m, w)
	}
	w.WriteMsg(m)
}


func main(){
	dns.HandleFunc(".", gelwormHandler)

	port := 15353
	server := &dns.Server{Addr: "0.0.0.0:" + strconv.Itoa(port), Net: "udp"}
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start: %s\n", err.Error())
	}
}

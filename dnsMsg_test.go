package dnsjson

import (
	"github.com/miekg/dns"
	"log"
	"testing"
)

func TestDNSMsg_DNS2JSON(t *testing.T) {

	msg := new(dns.Msg)
	msg.MsgHdr.Response = true
	msg.MsgHdr.RecursionAvailable = true
	msg.MsgHdr.AuthenticatedData = false
	msg.Answer = []dns.RR{
		CNAME("www.baidu.com 1097 CNAME www.a.shifen.com"),
		A("www.a.shifen.com. 62 A 220.181.38.149"),
		A("www.a.shifen.com. 62 A 220.181.38.150"),
	}

	dnsMsg := &DNSMsg{msg}

	b, _ := dnsMsg.DNS2JSON()
	log.Println(string(b))

}

func TestDNSMsg_JSON2DNS(t *testing.T) {

}

// A returns an A record from rr. It panics on errors.
func A(rr string) *dns.A { r, _ := dns.NewRR(rr); return r.(*dns.A) }

// AAAA returns an AAAA record from rr. It panics on errors.
func AAAA(rr string) *dns.AAAA { r, _ := dns.NewRR(rr); return r.(*dns.AAAA) }

// CNAME returns a CNAME record from rr. It panics on errors.
func CNAME(rr string) *dns.CNAME { r, _ := dns.NewRR(rr); return r.(*dns.CNAME) }

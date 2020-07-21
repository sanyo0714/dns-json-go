package dnsjson

import (
	"encoding/json"
	"github.com/miekg/dns"
	"github.com/sanyo0714/dns-json-go/jsondns"
	"log"
	"net"
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

	msg.SetQuestion("www.baidu.com", dns.TypeA)

	edns0Subnet := new(dns.EDNS0_SUBNET)
	edns0Subnet.Code = dns.EDNS0SUBNET
	edns0Subnet.SourceScope = 0
	edns0Subnet.Family = uint16(1)
	edns0Subnet.SourceNetmask = uint8(24)
	edns0Subnet.Address = net.ParseIP("10.10.10.10")

	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	opt.SetUDPSize(dns.DefaultMsgSize)
	opt.Option = append(opt.Option, edns0Subnet)

	msg.Extra = append(msg.Extra, opt)

	dnsMsg := &DNSMsg{msg}

	jsonByte, _ := dnsMsg.DNS2JSON()
	log.Println(string(jsonByte))

}

func TestDNSMsg_JSON2DNS(t *testing.T) {
	msg := new(dns.Msg)
	msg.MsgHdr.Response = true
	msg.MsgHdr.RecursionAvailable = true
	msg.MsgHdr.AuthenticatedData = false
	msg.Answer = []dns.RR{
		CNAME("www.baidu.com 1097 CNAME www.a.shifen.com"),
		A("www.a.shifen.com. 62 A 220.181.38.149"),
		A("www.a.shifen.com. 62 A 220.181.38.150"),
	}

	msg.SetQuestion("www.baidu.com", dns.TypeA)

	edns0Subnet := new(dns.EDNS0_SUBNET)
	edns0Subnet.Code = dns.EDNS0SUBNET
	edns0Subnet.SourceScope = 0
	edns0Subnet.Family = uint16(1)
	edns0Subnet.SourceNetmask = uint8(24)
	edns0Subnet.Address = net.ParseIP("10.10.10.10")

	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	opt.SetUDPSize(dns.DefaultMsgSize)
	opt.Option = append(opt.Option, edns0Subnet)

	msg.Extra = append(msg.Extra, opt)

	dnsMsg := &DNSMsg{msg}

	jsonByte, _ := dnsMsg.DNS2JSON()
	log.Println(string(jsonByte))
	// -------------------------------
	jsonMsg := &jsondns.JSONMsg{}
	json.Unmarshal(jsonByte, jsonMsg)

	log.Println(jsonMsg)

}

// A returns an A record from rr. It panics on errors.
func A(rr string) *dns.A { r, _ := dns.NewRR(rr); return r.(*dns.A) }

// AAAA returns an AAAA record from rr. It panics on errors.
func AAAA(rr string) *dns.AAAA { r, _ := dns.NewRR(rr); return r.(*dns.AAAA) }

// CNAME returns a CNAME record from rr. It panics on errors.
func CNAME(rr string) *dns.CNAME { r, _ := dns.NewRR(rr); return r.(*dns.CNAME) }

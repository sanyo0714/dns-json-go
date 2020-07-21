package dnsjson

import (
	"errors"
	"github.com/miekg/dns"
	"github.com/sanyo0714/dns-json-go/jsondns"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// DNSMsg is contains dns.Msg
type DNSMsg struct {
	*dns.Msg
}

// DNS2JSON is a function dns message to json
func (msg *DNSMsg) DNS2JSON() (*jsondns.JSONMsg, error) {

	now := time.Now().UTC()

	resp := new(jsondns.JSONMsg)
	resp.Status = uint32(msg.Rcode)
	resp.TC = msg.Truncated
	resp.RD = msg.RecursionDesired
	resp.RA = msg.RecursionAvailable
	resp.AD = msg.AuthenticatedData
	resp.CD = msg.CheckingDisabled

	resp.Question = make([]jsondns.Question, 0, len(msg.Question))
	for _, question := range msg.Question {
		jsonQuestion := jsondns.Question{
			Name: question.Name,
			Type: question.Qtype,
		}
		resp.Question = append(resp.Question, jsonQuestion)
	}

	resp.Answer = make([]jsondns.JSONRR, 0, len(msg.Answer))
	for _, rr := range msg.Answer {
		jsonAnswer := &jsondns.JSONRR{}
		jsonAnswer.MarshalRR(rr, now)
		if !resp.HaveTTL || jsonAnswer.TTL < resp.LeastTTL {
			resp.HaveTTL = true
			resp.LeastTTL = jsonAnswer.TTL
			resp.EarliestExpires = jsonAnswer.Expires
		}
		resp.Answer = append(resp.Answer, *jsonAnswer)
	}

	resp.Authority = make([]jsondns.JSONRR, 0, len(msg.Ns))
	for _, rr := range msg.Ns {
		jsonAuthority := &jsondns.JSONRR{}
		jsonAuthority.MarshalRR(rr, now)
		if !resp.HaveTTL || jsonAuthority.TTL < resp.LeastTTL {
			resp.HaveTTL = true
			resp.LeastTTL = jsonAuthority.TTL
			resp.EarliestExpires = jsonAuthority.Expires
		}
		resp.Authority = append(resp.Authority, *jsonAuthority)
	}

	resp.Additional = make([]jsondns.JSONRR, 0, len(msg.Extra))
	for _, rr := range msg.Extra {
		jsonAdditional := &jsondns.JSONRR{}
		jsonAdditional.MarshalRR(rr, now)
		header := rr.Header()
		if header.Rrtype == dns.TypeOPT {
			opt := rr.(*dns.OPT)
			resp.Status = ((opt.Hdr.Ttl & 0xff000000) >> 20) | (resp.Status & 0xff)
			for _, option := range opt.Option {
				if option.Option() == dns.EDNS0SUBNET {
					edns0 := option.(*dns.EDNS0_SUBNET)
					clientAddress := edns0.Address
					if clientAddress == nil {
						clientAddress = net.IP{0, 0, 0, 0}
					} else if ipv4 := clientAddress.To4(); ipv4 != nil {
						clientAddress = ipv4
					}
					resp.ECS = clientAddress.String() + "/" + strconv.FormatUint(uint64(edns0.SourceScope), 10)
				}
			}
			continue
		}
		if !resp.HaveTTL || jsonAdditional.TTL < resp.LeastTTL {
			resp.HaveTTL = true
			resp.LeastTTL = jsonAdditional.TTL
			resp.EarliestExpires = jsonAdditional.Expires
		}
		resp.Additional = append(resp.Additional, *jsonAdditional)
	}

	return resp, nil
}

// JSON2DNS is a function json to dns message
func (msg *DNSMsg) JSON2DNS(jsonMsg *jsondns.JSONMsg, udpSize uint16, ednsClientNetmask uint8) (reply *dns.Msg, err error) {

	now := time.Now().UTC()

	reply = msg.Copy()
	reply.Truncated = jsonMsg.TC
	reply.AuthenticatedData = jsonMsg.AD
	reply.CheckingDisabled = jsonMsg.CD
	reply.Rcode = dns.RcodeServerFailure

	reply.Answer = make([]dns.RR, 0, len(jsonMsg.Answer))
	for _, rr := range jsonMsg.Answer {
		dnsRR, err := rr.UnmarshalRR(now)
		if err != nil {
			log.Println(err)
		} else {
			reply.Answer = append(reply.Answer, dnsRR)
		}
	}

	reply.Ns = make([]dns.RR, 0, len(jsonMsg.Authority))
	for _, rr := range jsonMsg.Authority {
		dnsRR, err := rr.UnmarshalRR(now)
		if err != nil {
			log.Println(err)
		} else {
			reply.Ns = append(reply.Ns, dnsRR)
		}
	}

	reply.Extra = make([]dns.RR, 0, len(jsonMsg.Additional)+1)
	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	if udpSize >= 512 {
		opt.SetUDPSize(udpSize)
	} else {
		opt.SetUDPSize(512)
	}
	opt.SetDo(false)
	ednsClientSubnet := jsonMsg.ECS
	ednsClientFamily := uint16(1)
	ednsClientAddress := net.IP(nil)
	ednsClientScope := uint8(255)
	if ednsClientSubnet != "" {
		slash := strings.IndexByte(ednsClientSubnet, '/')
		if slash < 0 {
			err = errors.New("Invalid client subnet")
			return
		}
		ednsClientAddress = net.ParseIP(ednsClientSubnet[:slash])
		if ednsClientAddress == nil {
			err = errors.New("Invalid client subnet address")
			return
		}
		if ipv4 := ednsClientAddress.To4(); ipv4 != nil {
			ednsClientAddress = ipv4
		} else {
			ednsClientFamily = 2
		}
		scope, parseErr := strconv.ParseUint(ednsClientSubnet[slash+1:], 10, 8)
		if parseErr != nil {
			err = errors.New("Invalid client subnet address")
			return
		}
		ednsClientScope = uint8(scope)
	}

	if ednsClientAddress != nil {
		if ednsClientNetmask == 255 {
			if ednsClientFamily == 1 {
				ednsClientNetmask = 24
			} else {
				ednsClientNetmask = 56
			}
		}
		edns0Subnet := new(dns.EDNS0_SUBNET)
		edns0Subnet.Code = dns.EDNS0SUBNET
		edns0Subnet.Family = ednsClientFamily
		edns0Subnet.SourceNetmask = ednsClientNetmask
		edns0Subnet.SourceScope = ednsClientScope
		edns0Subnet.Address = ednsClientAddress
		opt.Option = append(opt.Option, edns0Subnet)
	}
	reply.Extra = append(reply.Extra, opt)
	for _, rr := range jsonMsg.Additional {
		dnsRR, err := rr.UnmarshalRR(now)
		if err != nil {
			log.Println(err)
		} else {
			reply.Extra = append(reply.Extra, dnsRR)
		}
	}

	reply.Rcode = int(jsonMsg.Status & 0xf)
	opt.Hdr.Ttl = (opt.Hdr.Ttl & 0x00ffffff) | ((jsonMsg.Status & 0xff0) << 20)
	reply.Extra[0] = opt

	return
}

package jsondns

import (
	"fmt"
	"github.com/miekg/dns"
	"strings"
)

// JSONMsg is dns json struct
type JSONMsg struct {
	// Standard DNS response code (32 bit integer).
	Status uint32 `json:"Status"`
	// Whether the response is truncated
	TC bool `json:"TC"`
	// Always true for Google Public DNS
	RD bool `json:"RD"`
	// Always true for Google Public DNS
	RA bool `json:"RA"`
	// Whether all response data was validated with DNSSEC
	AD bool `json:"AD"`
	// 	// Whether the client asked to disable DNSSEC
	CD         bool       `json:"CD"`
	Question   []Question `json:"Question"`
	Answer     []JSONRR   `json:"Answer,omitempty"`
	Authority  []JSONRR   `json:"Authority,omitempty"`
	Additional []JSONRR   `json:"Additional,omitempty"`
	ECS        string     `json:"edns_client_subnet,omitempty"`
	Comment    string     `json:"Comment,omitempty"`
}

// Question is dns question json struct
type Question struct {
	// FQDN with trailing dot
	Name string `json:"name"`
	// Standard DNS JSONRR type
	Type uint16 `json:"type"`
}

// JSONRR is dns rr json struct
type JSONRR struct {
	// Always matches name in the Question section
	Name string `json:"name"`
	// Standard DNS JSONRR type
	Type uint16 `json:"type"`
	// Record's time-to-live in seconds
	TTL uint32 `json:"TTL"`
	// Data for A - IP address as text
	Data string `json:"data"`
}

// MarshalRR is marshal dns RR
func (jsonRR *JSONRR) MarshalRR(rr dns.RR) {
	rrHeader := rr.Header()
	jsonRR.Name = rrHeader.Name
	jsonRR.Type = rrHeader.Rrtype
	jsonRR.TTL = rrHeader.Ttl
	data := strings.SplitN(rr.String(), "\t", 5)
	if len(data) >= 5 {
		jsonRR.Data = data[4]
	}
}

// UnmarshalRR is unmarshal  dns RR
func (jsonRR *JSONRR) UnmarshalRR() (dnsRR dns.RR, err error) {
	if strings.ContainsAny(jsonRR.Name, "\t\r\n \"();\\") {
		return nil, UnmarshalError{fmt.Sprintf("Record name contains space: %q", jsonRR.Name)}
	}
	rrType, ok := dns.TypeToString[jsonRR.Type]
	if !ok {
		return nil, UnmarshalError{fmt.Sprintf("Unknown record type: %d", jsonRR.Type)}
	}
	if strings.ContainsAny(jsonRR.Data, "\r\n") {
		return nil, UnmarshalError{fmt.Sprintf("Record data contains newline: %q", jsonRR.Data)}
	}
	zone := fmt.Sprintf("%s %d IN %s %s", jsonRR.Name, jsonRR.TTL, rrType, jsonRR.Data)
	dnsRR, err = dns.NewRR(zone)
	return
}

// UnmarshalError is unmarshal error message
type UnmarshalError struct {
	err string
}

// Error is ...
func (e UnmarshalError) Error() string {
	return "json-dns: " + e.err
}

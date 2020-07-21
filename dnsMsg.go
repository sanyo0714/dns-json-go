package dnsjson

import (
	"encoding/json"
	"github.com/miekg/dns"
	"github.com/sanyo0714/dns-json-go/jsondns"
)

// DNSMsg is contains dns.Msg
type DNSMsg struct {
	dns.Msg
}

// DNS2JSON is a function dns message to json
func (msg *DNSMsg) DNS2JSON() ([]byte, error) {
	jsonMsg := new(jsondns.JSONMsg)

	return json.Marshal(jsonMsg)
}

// JSON2DNS is a function json to dns message
func (msg *DNSMsg) JSON2DNS(jsonMsg *jsondns.JSONMsg) {

}

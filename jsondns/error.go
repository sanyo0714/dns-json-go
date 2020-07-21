package jsondns

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/miekg/dns"
)

type dnsError struct {
	Status  uint32 `json:"Status"`
	Comment string `json:"Comment,omitempty"`
}

func FormatError(w http.ResponseWriter, comment string, errcode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	errJson := dnsError{
		Status:  dns.RcodeServerFailure,
		Comment: comment,
	}
	errStr, err := json.Marshal(errJson)
	if err != nil {
		log.Fatalln(err)
	}
	w.WriteHeader(errcode)
	w.Write(errStr)
}

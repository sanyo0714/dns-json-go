package jsondns

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
	Answer     []RR       `json:"Answer,omitempty"`
	Authority  []RR       `json:"Authority,omitempty"`
	Additional []RR       `json:"Additional,omitempty"`
	ECS        string     `json:"edns_client_subnet,omitempty"`
	Comment    string     `json:"Comment,omitempty"`
}

// Question is dns question json struct
type Question struct {
	// FQDN with trailing dot
	Name string `json:"name"`
	// Standard DNS RR type
	Type uint16 `json:"type"`
}

// RR is dns rr json struct
type RR struct {
	// Always matches name in the Question section
	Name string `json:"name"`
	// Standard DNS RR type
	Type uint16 `json:"type"`
	// Record's time-to-live in seconds
	TTL uint32 `json:"TTL"`
	// Data for A - IP address as text
	Data string `json:"data"`
}

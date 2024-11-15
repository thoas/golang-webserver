package main

type request struct {
	URI      string            `json:"uri"`
	ClientIP string            `json:"client_ip,omitempty"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers,omitempty"`
	Json     interface{}       `json:"json,omitempty"`
	Body     string            `json:"body,omitempty"`
}

type response struct {
	Requests []*request `json:"requests"`
}

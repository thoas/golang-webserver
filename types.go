package main

import "time"

type Request struct {
	URI       string            `json:"uri"`
	ClientIP  string            `json:"client_ip,omitempty"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	Json      interface{}       `json:"json,omitempty"`
	Body      string            `json:"body,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type Response struct {
	Requests []*Request `json:"requests"`
}

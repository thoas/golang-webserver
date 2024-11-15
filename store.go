package main

import (
	"io"
	"net/http"
	"sync"
)

type store struct {
	mu       sync.Mutex
	requests []*request
}

func newStore() *store {
	return &store{
		mu:       sync.Mutex{},
		requests: make([]*request, 0),
	}
}

func (s *store) add(r *http.Request) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := &request{
		Method:   r.Method,
		URI:      r.RequestURI,
		Headers:  make(map[string]string),
		ClientIP: remoteIP(r),
	}

	for key, vals := range r.Header {
		req.Headers[key] = vals[0]
	}

	var (
		err  error
		save io.ReadCloser
	)
	save, r.Body, err = drainBody(r.Body)
	if err != nil {
		return err
	}
	if save != nil {
		if r.Header.Get("Content-Type") == "application/json" {
			out, err := dumpJson(save)
			if err != nil {
				return err
			}
			req.Json = out
		} else {
			body, err := io.ReadAll(save)
			if err != nil {
				return err
			}
			req.Body = string(body)
		}
	}

	s.requests = append(s.requests, req)

	return nil
}

func (s *store) flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = make([]*request, 0)
}

func (s *store) dump() []*request {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.requests
}

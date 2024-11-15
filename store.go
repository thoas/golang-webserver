package main

import (
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type store struct {
	mu       sync.Mutex
	requests []*request
	capacity int
	index    int
}

func newStore(capacity int) *store {
	return &store{
		mu:       sync.Mutex{},
		requests: make([]*request, capacity, capacity),
		capacity: capacity,
		index:    0,
	}
}

func (s *store) add(r *http.Request) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := &request{
		Method:    r.Method,
		URI:       r.RequestURI,
		Headers:   make(map[string]string),
		ClientIP:  remoteIP(r),
		CreatedAt: time.Now(),
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

	s.requests[s.index] = req
	s.index++
	if s.index >= s.capacity {
		s.index = 0
	}

	return nil
}

func (s *store) flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = make([]*request, 0, s.capacity)
}

func (s *store) dump() []*request {
	s.mu.Lock()
	requests := make([]*request, 0)
	for i := range s.requests {
		if s.requests[i] == nil {
			continue
		}
		requests = append(requests, s.requests[i])
	}
	s.mu.Unlock()

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].CreatedAt.Before(requests[j].CreatedAt)
	})

	return requests
}

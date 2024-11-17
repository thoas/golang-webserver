package main

import (
	"net/http"
	"slices"
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu              sync.Mutex
	requests        []*Request
	capacity        int
	index           int
	excludedHeaders []string
}

func NewStore(capacity int, excludedHeaders []string) *Store {
	return &Store{
		mu:              sync.Mutex{},
		requests:        make([]*Request, capacity, capacity),
		capacity:        capacity,
		index:           0,
		excludedHeaders: excludedHeaders,
	}
}

func (s *Store) Add(r *http.Request) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req := &Request{
		Method:    r.Method,
		URI:       r.RequestURI,
		Headers:   make(map[string]string),
		ClientIP:  remoteIP(r),
		CreatedAt: time.Now(),
	}

	for key, vals := range r.Header {
		if slices.Contains(s.excludedHeaders, key) {
			continue
		}
		req.Headers[key] = vals[0]
	}

	out, body, err := getBody(r)
	if err != nil {
		return err
	}

	if body != nil {
		req.Body = string(body)
	}
	if out != nil {
		req.Json = out
	}

	s.requests[s.index] = req
	s.index++
	if s.index >= s.capacity {
		s.index = 0
	}

	return nil
}

func (s *Store) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = make([]*Request, 0, s.capacity)
}

func (s *Store) Dump() []*Request {
	s.mu.Lock()
	requests := make([]*Request, 0)
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

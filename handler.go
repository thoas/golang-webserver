package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type handler struct {
	store *store
}

func (h *handler) logError(err interface{}) {
	log.Print(err)
}

func (h *handler) wrap(wrapped http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				h.logError(r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		wrapped(w, r)
	}
}

func (h *handler) root(w http.ResponseWriter, r *http.Request) {
	if err := h.store.add(r); err != nil {
		h.logError(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) flush(w http.ResponseWriter, _ *http.Request) {
	h.store.flush()
	w.WriteHeader(http.StatusOK)
}

func (h *handler) dump(w http.ResponseWriter, r *http.Request) {
	requests := h.store.dump()
	res := &response{Requests: requests}
	out, err := json.Marshal(res)
	if err != nil {
		h.logError(err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(out))
}

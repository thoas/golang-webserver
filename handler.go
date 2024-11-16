package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type handler struct {
	store *store
}

func (h *handler) logError(err interface{}) {
	log.Print(err)
}

func (h *handler) wrap(wrapped HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				h.logError(r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		if err := wrapped(w, r); err != nil {
			h.logError(err)
		}
	}
}

func (h *handler) root(w http.ResponseWriter, r *http.Request) error {
	if err := h.store.add(r); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) flush(w http.ResponseWriter, _ *http.Request) error {
	h.store.flush()
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) dump(w http.ResponseWriter, r *http.Request) error {
	requests := h.store.dump()
	res := &response{Requests: requests}
	out, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(out))

	return nil
}

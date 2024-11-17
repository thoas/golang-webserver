package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

type Handler struct {
	store *Store
}

func (h *Handler) logError(err interface{}) {
	log.Print(err)
}

func (h *Handler) Wrap(wrapped HandlerFunc) http.HandlerFunc {
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

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) error {
	if err := h.store.Add(r); err != nil {
		return err
	}

	fmt.Fprint(w, "Ok")

	return nil
}

func (h *Handler) Flush(w http.ResponseWriter, _ *http.Request) error {
	h.store.Flush()
	fmt.Fprint(w, "Ok")

	return nil
}

func (h *Handler) Dump(w http.ResponseWriter, r *http.Request) error {
	requests := h.store.Dump()
	res := &Response{Requests: requests}
	out, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(out))

	return nil
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DEFAULT_PORT     = 8080
	DEFAULT_CAPACITY = 100
)

func main() {
	var (
		port      int = DEFAULT_PORT
		capacity  int = DEFAULT_CAPACITY
		err       error
		sport     = os.Getenv("PORT")
		scapacity = os.Getenv("CAPACITY")
	)

	if sport != "" {
		port, err = strconv.Atoi(sport)

		if err != nil {
			log.Fatal(err)
		}
	}

	if scapacity != "" {
		capacity, err = strconv.Atoi(scapacity)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Running on :%d", port)

	h := NewHandler(NewStore(capacity))
	http.HandleFunc("/", h.Wrap(h.Root))
	http.HandleFunc("/dump", h.Wrap(h.Dump))
	http.HandleFunc("/flush", h.Wrap(h.Flush))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

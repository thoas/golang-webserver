package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const DEFAULT_PORT = 8080
const DEFAULT_CAPACITY = 100

func main() {
	var (
		port     int = DEFAULT_PORT
		capacity int = DEFAULT_CAPACITY
		err      error
	)

	sport := os.Getenv("PORT")
	scapacity := os.Getenv("CAPACITY")

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

	h := handler{store: newStore(capacity)}
	http.HandleFunc("/", h.wrap(h.root))
	http.HandleFunc("/dump", h.wrap(h.dump))
	http.HandleFunc("/flush", h.wrap(h.flush))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const DEFAULT_PORT = 8080

func main() {
	var port int = DEFAULT_PORT
	var err error

	sport := os.Getenv("PORT")

	if sport != "" {
		port, err = strconv.Atoi(sport)

		if err != nil {
			panic(err)
		}
	}

	log.Printf("Running on :%d", port)

	h := handler{store: newStore()}
	http.HandleFunc("/", h.root)
	http.HandleFunc("/dump", h.dump)
	http.HandleFunc("/flush", h.flush)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

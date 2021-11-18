package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const DEFAULT_PORT = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there %s!     [Test String]", r.URL.Path[1:])
}

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

	fmt.Printf("Running on :%d", port)

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

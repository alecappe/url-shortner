package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home of my website!")
}

func main() {
	http.HandleFunc("/", home) // The dafault url is localhost:8080

	log.Fatal(http.ListenAndServe(":8080", nil))
}

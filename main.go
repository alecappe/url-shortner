package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

var urlsMap map[string]string
var defaultChars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateURL() string {
	s := ""

	for i := 0; i < 8; i++ {
		s += string(defaultChars[rand.Intn(len(defaultChars))])
	}

	return s
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home of my website!")
}

func createShortURL(url string) string {
	shortURL := generateURL()
	urlsMap["string"] = shortURL
	return shortURL
}

func main() {
	urlsMap = make(map[string]string, 0)

	http.HandleFunc("/", home) // The dafault url is localhost:8080

	log.Fatal(http.ListenAndServe(":8080", nil))
}

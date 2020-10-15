package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

func createShortURL(url string) string {
	shortURL := "/" + generateURL()
	urlsMap[shortURL] = url
	return shortURL
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := string(r.URL.Path)

	url = strings.Replace(url, "/shorten/", "", 1)
	fmt.Println("Stripped url is:", url)

	shortURL := createShortURL(url)
	fmt.Fprintf(w, shortURL)

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home of my website!\n\n")
	url := string(r.URL.Path)
	expandedURL := urlsMap[url]
	if expandedURL != "" {
		fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
	}
}

func main() {
	urlsMap = make(map[string]string, 0)

	http.HandleFunc("/", home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

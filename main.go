package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

type urlsStruct struct {
	urls           map[string]string
	mux            sync.Mutex
}

var defaultChars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateURL() string {
	s := ""

	for i := 0; i < 8; i++ {
		s += string(defaultChars[rand.Intn(len(defaultChars))])
	}

	return s
}

func (u *urlsStruct) createShortURL(url string) string {
	shortURL := "/" + generateURL()
	u.mux.Lock()
	u.urls[shortURL] = url
	u.mux.Unlock()
	return shortURL
}

func (u *urlsStruct) handler(w http.ResponseWriter, r *http.Request) {
	url := string(r.URL.Path)

	url = strings.Replace(url, "/shorten/", "", 1)
	fmt.Println("Stripped url is:", url)

	shortURL := u.createShortURL(url)
	fmt.Fprintf(w, shortURL)

}

func (u *urlsStruct) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home of my website!\n\n")
	url := string(r.URL.Path)
	expandedURL := u.urls[url]
	if expandedURL != "" {
		fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
	}
}

func main() {
	data := urlsStruct{}

	data.urls = make(map[string]string, 0)
	http.HandleFunc("/", data.home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", data.handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
	homeCallCount  int
	shortCallCount int
	statsCallCount int
	urlsCount      int
	succRedir      int
	failRedir      int
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
	u.urlsCount++
	u.mux.Unlock()
	return shortURL
}

func (u *urlsStruct) handler(w http.ResponseWriter, r *http.Request) {
	u.shortCallCount++

	url := string(r.URL.Path)

	url = strings.Replace(url, "/shorten/", "", 1)
	fmt.Println("Stripped url is:", url)

	shortURL := u.createShortURL(url)
	fmt.Fprintf(w, shortURL)

}

func (u *urlsStruct) stats(w http.ResponseWriter, r *http.Request) {
	u.statsCallCount++

	fmt.Fprintf(w, "home called: %d\n", u.homeCallCount)
	fmt.Fprintf(w, "Shorten called: %d\n", u.shortCallCount)
	fmt.Fprintf(w, "Stats called: %d\n", u.statsCallCount)

	fmt.Fprintf(w, "Generated urls: %d\n", u.urlsCount)
	fmt.Fprintf(w, "Success redirect: %d\n", u.succRedir)
	fmt.Fprintf(w, "Failed redirect: %d\n", u.failRedir)
}

func (u *urlsStruct) home(w http.ResponseWriter, r *http.Request) {
	u.homeCallCount++

	fmt.Fprintf(w, "This is the home of my website!\n\n")

	url := string(r.URL.Path)
	if url != "/" {
		expandedURL := u.urls[url]
		if expandedURL != "" {
			fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
			u.succRedir++
			return
		}

		u.failRedir++
	}
}

func main() {
	data := urlsStruct{}

	data.urls = make(map[string]string, 0)

	// Init stats
	data.shortCallCount = 0
	data.homeCallCount = 0
	data.statsCallCount = 0
	data.urlsCount = 0
	data.succRedir = 0
	data.failRedir = 0

	// API
	http.HandleFunc("/", data.home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", data.handler)
	http.HandleFunc("/stats", data.stats)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

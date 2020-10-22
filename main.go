package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

type urlsStruct struct {
	urls map[string]string
	mux  sync.Mutex

	Stats struct {
		HomeVisit       int
		ShortenCall     int
		StatsVisit      int
		UrlsGenerated   int
		SuccessRedirect int
		FailedRedirect  int
	}
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
	u.Stats.UrlsGenerated++
	u.mux.Unlock()
	return shortURL
}

func (u *urlsStruct) handler(w http.ResponseWriter, r *http.Request) {
	u.Stats.ShortenCall++

	url := string(r.URL.Path)

	url = strings.Replace(url, "/shorten/", "", 1)
	fmt.Println("Stripped url is:", url)

	shortURL := u.createShortURL(url)
	fmt.Fprintf(w, shortURL)

}

func (u *urlsStruct) stats(w http.ResponseWriter, r *http.Request) {
	u.Stats.StatsVisit++

	formatNeeded, ok := r.URL.Query()["format"]
	if ok && formatNeeded[0] == "json" {
		b, err := json.Marshal(u)
		if err != nil {
			log.Fatalf("Unable to encode")
		}
		fmt.Fprintf(w, string(b))
		return
	}

	fmt.Fprintf(w, "home called: %d\n", u.Stats.HomeVisit)
	fmt.Fprintf(w, "Shorten called: %d\n", u.Stats.ShortenCall)
	fmt.Fprintf(w, "Stats called: %d\n", u.Stats.StatsVisit)

	fmt.Fprintf(w, "Generated urls: %d\n", u.Stats.UrlsGenerated)
	fmt.Fprintf(w, "Success redirect: %d\n", u.Stats.SuccessRedirect)
	fmt.Fprintf(w, "Failed redirect: %d\n", u.Stats.FailedRedirect)
}

func (u *urlsStruct) home(w http.ResponseWriter, r *http.Request) {
	u.Stats.HomeVisit++

	fmt.Fprintf(w, "This is the home of my website!\n\n")

	url := string(r.URL.Path)
	if url != "/" {
		expandedURL := u.urls[url]
		if expandedURL != "" {
			fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
			u.Stats.SuccessRedirect++
			return
		}

		u.Stats.FailedRedirect++
	}
}

func main() {
	data := urlsStruct{}

	data.urls = make(map[string]string, 0)

	// Init stats
	data.Stats.ShortenCall = 0
	data.Stats.HomeVisit = 0
	data.Stats.StatsVisit = 0
	data.Stats.UrlsGenerated = 0
	data.Stats.SuccessRedirect = 0
	data.Stats.FailedRedirect = 0

	// API
	http.HandleFunc("/", data.home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", data.handler)
	http.HandleFunc("/stats", data.stats)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

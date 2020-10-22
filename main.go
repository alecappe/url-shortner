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
	mux  sync.RWMutex

	Stats struct {
		HomeVisit       int
		ShortenCall     int
		StatsVisit      int
		UrlsGenerated   int
		SuccessRedirect int
		FailedRedirect  int
	}
}

func newUrlsStruct() *urlsStruct {
	v := urlsStruct{}
	v.urls = make(map[string]string)
	v.Stats.ShortenCall = 0
	v.Stats.HomeVisit = 0
	v.Stats.StatsVisit = 0
	v.Stats.UrlsGenerated = 0
	v.Stats.SuccessRedirect = 0
	v.Stats.FailedRedirect = 0
	return &v
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
	defer u.mux.Unlock()
	u.urls[shortURL] = url
	u.Stats.UrlsGenerated++
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
		u.mux.RLock()
		b, err := json.Marshal(u)
		u.mux.RUnlock()
		if err != nil {
			log.Fatalf("Unable to encode")
		}
		fmt.Fprintf(w, string(b))
		return
	}

	u.mux.RLock()
	fmt.Fprintf(w, "home called: %d\n", u.Stats.HomeVisit)
	fmt.Fprintf(w, "Shorten called: %d\n", u.Stats.ShortenCall)
	fmt.Fprintf(w, "Stats called: %d\n", u.Stats.StatsVisit)

	fmt.Fprintf(w, "Generated urls: %d\n", u.Stats.UrlsGenerated)
	fmt.Fprintf(w, "Success redirect: %d\n", u.Stats.SuccessRedirect)
	fmt.Fprintf(w, "Failed redirect: %d\n", u.Stats.FailedRedirect)
	u.mux.RUnlock()
}

func (u *urlsStruct) home(w http.ResponseWriter, r *http.Request) {
	u.Stats.HomeVisit++

	fmt.Fprintf(w, "This is the home of my website!\n\n")

	url := string(r.URL.Path)
	if url != "/" {
		u.mux.RLock()
		expandedURL := u.urls[url]
		u.mux.RUnlock()
		if expandedURL != "" {
			fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
			u.Stats.SuccessRedirect++
			return
		}

		u.Stats.FailedRedirect++
	}
}

func main() {
	data := newUrlsStruct()

	// API
	http.HandleFunc("/", data.home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", data.handler)
	http.HandleFunc("/stats", data.stats)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

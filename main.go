package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

type urlsStruct struct {
	urls map[string]string
	mux  sync.RWMutex

	Stats struct {
		HomeVisit       int32
		ShortenCall     int32
		StatsVisit      int32
		UrlsGenerated   int32
		SuccessRedirect int32
		FailedRedirect  int32
	}
}

func newUrlsStruct() *urlsStruct {
	v := urlsStruct{}
	v.urls = make(map[string]string)
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
	atomic.AddInt32(&u.Stats.UrlsGenerated, 1)
	return shortURL
}

func (u *urlsStruct) handler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&u.Stats.ShortenCall, 1)

	url := string(r.URL.Path)

	url = strings.Replace(url, "/shorten/", "", 1)
	fmt.Println("Stripped url is:", url)

	shortURL := u.createShortURL(url)
	fmt.Fprintf(w, shortURL)
}

func (u *urlsStruct) showStats(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&u.Stats.StatsVisit, 1)

	formatNeeded, ok := r.URL.Query()["format"]
	if ok && formatNeeded[0] == "json" {
		u.mux.RLock()
		b, err := json.MarshalIndent(u, "", "    ")
		u.mux.RUnlock()
		if err != nil {
			log.Fatalf("Unable to encode")
		}
		fmt.Fprintf(w, string(b))
		return
	}

	u.mux.RLock()
	fmt.Fprintf(w, "Home called: %d\n", atomic.LoadInt32(&u.Stats.HomeVisit))
	fmt.Fprintf(w, "Shorten called: %d\n", atomic.LoadInt32(&u.Stats.ShortenCall))
	fmt.Fprintf(w, "Stats called: %d\n", atomic.LoadInt32(&u.Stats.StatsVisit))

	fmt.Fprintf(w, "Generated urls: %d\n", atomic.LoadInt32(&u.Stats.UrlsGenerated))
	fmt.Fprintf(w, "Success redirect: %d\n", atomic.LoadInt32(&u.Stats.SuccessRedirect))
	fmt.Fprintf(w, "Failed redirect: %d\n", atomic.LoadInt32(&u.Stats.FailedRedirect))
	u.mux.RUnlock()
}

func (u *urlsStruct) home(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&u.Stats.HomeVisit, 1)

	fmt.Fprintf(w, "This is the home of my website!\n\n")

	url := string(r.URL.Path)
	if url != "/" {
		u.mux.RLock()
		expandedURL := u.urls[url]
		u.mux.RUnlock()
		if expandedURL != "" {
			fmt.Fprintf(w, "Redirect to:\n"+expandedURL)
			atomic.AddInt32(&u.Stats.SuccessRedirect, 1)
			return
		}

		atomic.AddInt32(&u.Stats.FailedRedirect, 1)
	}
}

func main() {
	var serverPort string
	flag.StringVar(&serverPort, "addr", "localhost:8080", "Use to set the server address")
	flag.Parse()

	data := newUrlsStruct()

	// API
	http.HandleFunc("/", data.home) // The dafault url is localhost:8080
	http.HandleFunc("/shorten/", data.handler)
	http.HandleFunc("/stats", data.showStats)

	log.Fatal(http.ListenAndServe(serverPort, nil))
}

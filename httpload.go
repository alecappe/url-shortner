package main

import (
	"flag"
	"fmt"
	"net/http"
)

type stats struct {
	failedRequest  int
	successRequest int
}

func (s *stats) req(url string) {
	r, err := http.Get(url)
	if err != nil {
		s.failedRequest++
		return
	}
	fmt.Println(r.StatusCode)

	s.successRequest++

}

func main() {
	URL := ""

	flag.StringVar(&URL, "", "http://www.google.it", "url to test")
	flag.Parse()

	s := stats{}
	s.successRequest = 0
	s.failedRequest = 0

	for i := 0; i < 10; i++ {
		go s.req(URL)
	}

	var input string
	fmt.Scanln(&input)
	fmt.Println(s.successRequest, s.failedRequest)
}

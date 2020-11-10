package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
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

func (s *stats) reqXn(url string, reqNum int) {
	for i := 0; i < reqNum; i++ {
		s.req(url)
	}
}

func (s *stats) reqXtime(url string, duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		s.failedRequest++
		return
	}

	stop := make(chan bool)

	go func() {
		for {
			s.req(url)
			select {
			case <-time.After(d):
			case <-stop:
				return
			}
		}
	}()
}

func main() {
	URL := ""
	concNum := 0   // number of concurrent workers
	reqNum := 0    // number of total requests
	duration := "" // duration of application to send requests

	flag.StringVar(&URL, "", "http://www.google.it", "url to test")
	flag.IntVar(&concNum, "w", 50, "number of workers to run concurrently. default:50")
	flag.IntVar(&reqNum, "n", 200, "number of requests to run. default:200.")
	flag.StringVar(&duration, "z", "", "duration of application to send requests.")
	flag.Parse()

	if reqNum < concNum {
		fmt.Println("The number of request can't be smaller than number of workers")
		os.Exit(1)
	}

	s := stats{}
	s.successRequest = 0
	s.failedRequest = 0

	for i := 0; i < concNum; i++ {
		if duration != "" {
			go s.reqXtime(URL, duration)
		} else {
			go s.reqXn(URL, reqNum/concNum)
		}
	}

	var input string
	fmt.Scanln(&input)
	fmt.Println(s.successRequest, s.failedRequest)
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type stats struct {
	failedRequest  int
	successRequest int
}

func (s *stats) req(url string, reqNum int) {
	for i := 0; i < reqNum; i++ {
		r, err := http.Get(url)
		if err != nil {
			s.failedRequest++
			return
		}
		fmt.Println(r.StatusCode)

		s.successRequest++
	}
}

func main() {
	URL := ""
	concNum := 0 // number of concurrent workers
	reqNum := 0  // number of total requests

	flag.StringVar(&URL, "", "http://www.google.it", "url to test")
	flag.IntVar(&concNum, "w", 50, "number of workers to run concurrently. default:50")
	flag.IntVar(&reqNum, "n", 200, "number of requests to run. default:200.")
	flag.Parse()

	if reqNum < concNum {
		fmt.Println("The number of request can't be smaller than number of workers")
		os.Exit(1)
	}

	s := stats{}
	s.successRequest = 0
	s.failedRequest = 0

	for i := 0; i < concNum; i++ {
		go s.req(URL, reqNum/concNum)
	}

	var input string
	fmt.Scanln(&input)
	fmt.Println(s.successRequest, s.failedRequest)
}

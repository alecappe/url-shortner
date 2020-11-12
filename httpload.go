package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type stats struct {
	failedRequest  int
	successRequest int

}

func (s *stats) reqXn(url string, reqNum int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < reqNum; i++ {
		r, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		fmt.Println(r.StatusCode)
		s.successRequest++
		r.Body.Close()
	}
}

func (s *stats) reqXtime(cancelCtx context.Context, url string) {
	client := http.Client{}

	for {
		r, err := http.NewRequestWithContext(cancelCtx, "GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		resp, err := client.Do(r)
		if err != nil {
			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		fmt.Println(resp.StatusCode)
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		s.successRequest++
	}
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

	if duration == "" && reqNum < concNum {
		fmt.Println("The number of requests can't be smaller than number of workers")
		os.Exit(1)
	}

	var d time.Duration
	var err error

	if duration != "" {
		d, err = time.ParseDuration(duration)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	s := stats{}

	wg := new(sync.WaitGroup)
	ctx := context.Background()

	for i := 0; i < concNum; i++ {
		if duration != "" {
			cancelCtx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			s.reqXtime(cancelCtx, URL)
		} else {
			wg.Add(1)
			go s.reqXn(URL, reqNum/concNum, wg)
		}
	}

	wg.Wait()

	fmt.Println(s.successRequest, s.failedRequest)
}

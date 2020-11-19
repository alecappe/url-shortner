package main

import (
	"context"
	"errors"
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

	reqTimes     []time.Duration
	totalTime    time.Duration
	slowest      time.Duration
	fastest      time.Duration
	average      time.Duration
	requestsXSec float64
}

func (s *stats) reqXn(url string, reqNum int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < reqNum; i++ {
		start := time.Now()
		r, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		r.Body.Close()
		s.successRequest++
		s.elapsedTime(start)
	}
}

func (s *stats) reqXtime(cancelCtx context.Context, url string) {
	client := http.Client{}

	for {
		start := time.Now()
		r, err := http.NewRequestWithContext(cancelCtx, "GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		resp, err := client.Do(r)
		if err != nil {
			err = errors.Unwrap(err)
			if err == context.DeadlineExceeded {
				return
			}

			fmt.Println(err.Error())
			s.failedRequest++
			return
		}

		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		s.successRequest++
		s.elapsedTime(start)
	}
}

func (s *stats) elapsedTime(start time.Time) {
	elapsed := time.Since(start)
	s.reqTimes = append(s.reqTimes, elapsed)

	if s.slowest == 0 && s.fastest == 0 {
		s.slowest, s.fastest = elapsed, elapsed
		return
	}

	if elapsed > s.slowest {
		s.slowest = elapsed
		return
	}

	if elapsed < s.fastest {
		s.fastest = elapsed
		return
	}
}

func (s *stats) averageTime() {
	var total time.Duration
	for _, d := range s.reqTimes {
		total = total + d
	}
	s.average = total / time.Duration(len(s.reqTimes))
}

func (s *stats) printStats() {
	// Calculate average before print stats
	s.averageTime()

	// Calculate requests/sec
	s.requestsXSec = float64(len(s.reqTimes)) / s.totalTime.Seconds()

	// Print stats
	fmt.Println("Summary:")
	fmt.Println("Total:", s.totalTime)
	fmt.Println("Slowest:", s.slowest)
	fmt.Println("Fastest:", s.fastest)
	fmt.Println("Average:", s.average)
	fmt.Println("Requests/sec:", s.requestsXSec)
}

func main() {
	URL := "https://www.google.it"
	concNum := 0   // number of concurrent workers
	reqNum := 0    // number of total requests
	duration := "" // duration of application to send requests

	// flag.StringVar(&URL, "", "http://www.google.it", "url to test")
	flag.IntVar(&concNum, "w", 50, "number of workers to run concurrently. default:50")
	flag.IntVar(&reqNum, "n", 200, "number of requests to run. default:200.")
	flag.StringVar(&duration, "z", "", "duration of application to send requests.")
	flag.Parse()

	if flag.Arg(0) != "" {
		URL = flag.Arg(0)
	}

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

	start := time.Now()

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

	s.totalTime = time.Since(start)
	s.printStats()
	fmt.Println(s.successRequest, s.failedRequest)
}

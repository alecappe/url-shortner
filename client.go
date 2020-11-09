package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type dataURLS struct {
	Stats struct {
		HomeVisit       int32
		ShortenCall     int32
		StatsVisit      int32
		UrlsGenerated   int32
		SuccessRedirect int32
		FailedRedirect  int32
	}
}

func startServer() *exec.Cmd {
	cmd := exec.Command("go", "run", "main.go", "-load", "urls.json")

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

func doTest() error {
	// Start the http server
	cmd := startServer()
	shouldKill := true

	defer func() {
		if shouldKill == true {
			cmd.Process.Signal(os.Interrupt)
		}
	}()

	time.Sleep(1000 * time.Millisecond)
	// +++++++++++++++++++++++++++++++++++++++

	// call / on an non-existing url and check the http.StatusCode
	client := &http.Client{}

	nonExistingURL := "http://localhost:8080/not-existing-url"

	resp, err := client.Get(nonExistingURL)
	if err != nil {
		return fmt.Errorf("Get request error", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf(nonExistingURL, "should not be found")
	}
	// +++++++++++++++++++++++++++++++++++++++

	// call / on a existing url (loaded from urls.json file at startup)
	d, err := ioutil.ReadFile("urls.json")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	var data map[string]string
	err = json.Unmarshal(d, &data)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	existingURL := "http://localhost:8080/urlnum1"
	resp, err = client.Get(existingURL)
	if err != nil {
		return fmt.Errorf("Get request error", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(existingURL, "should be found")
	}
	// +++++++++++++++++++++++++++++++++++++++

	// call /shorten with a new URL (that wasn't in urls.json)
	shortenURL := "http://localhost:8080/shorten/test.com"
	resp, err = client.Get(shortenURL)
	if err != nil {
		return fmt.Errorf("Get request error", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(shortenURL, "didn't create a url")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	generatedURL := string(body)

	// call / with an URL that has just been added (check status code and ensure the redirect is working)
	testGeneratedURL := "http://localhost:8080" + generatedURL
	resp2, err := client.Get(testGeneratedURL)
	if err != nil {
		return fmt.Errorf("Get request error", err)
	}
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(testGeneratedURL, "don't exists")
	}
	// +++++++++++++++++++++++++++++++++++++++

	// call /statistics, unmarshall the returned json and checks that it corresponds to the actions taken
	statisticsURL := "http://localhost:8080/stats?format=json"
	resp2, err = client.Get(statisticsURL)
	if err != nil {
		return fmt.Errorf("Get request error", err)
	}
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(statisticsURL, "not found")
	}

	defer resp2.Body.Close()
	body, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("error")
	}

	s := dataURLS{}
	json.Unmarshal([]byte(body), &s)

	if s.Stats.HomeVisit != 0 &&
		s.Stats.ShortenCall != 1 &&
		s.Stats.StatsVisit != 1 &&
		s.Stats.UrlsGenerated != 5 &&
		s.Stats.SuccessRedirect != 2 &&
		s.Stats.FailedRedirect != 1 {
		return fmt.Errorf("Unexpected statistics")
	}

	// +++++++++++++++++++++++++++++++++++++++

	// terminate the http server with a signal
	fmt.Println("Client finish without errors")
	cmd.Process.Signal(os.Interrupt)
	shouldKill = false

	return nil
	// +++++++++++++++++++++++++++++++++++++++
}

func main() {
	err := doTest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

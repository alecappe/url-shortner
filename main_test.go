package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateURL(t *testing.T) {
	v := generateURL()
	if len(v) != 8 {
		t.Errorf("The generated url have a unexpected length")
	}
}

func TestCreateShortUrl(t *testing.T) {
	d := newUrlsStruct()

	url := d.createShortURL("www.test.com")
	if len(url) != 9 {
		t.Errorf("The generated url have a unexpected length")
	}

	if d.urls[url] != "www.test.com" {
		t.Errorf("The saved url in the map is wrong")
	}
}

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/shorten/www.google.it", nil)
	if err != nil {
		t.Fatal(err)
	}

	data := newUrlsStruct()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(data.handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if len(rr.Body.String()) != 9 {
		t.Errorf("Handler returned unexpected string length")
	}

	if data.urls[rr.Body.String()] != "www.google.it" {
		t.Errorf("The saved url in the map is wrong")
	}
}

func TestShowStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	data := newUrlsStruct()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(data.showStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Home called: 0\nShorten called: 0\nStats called: 1\nGenerated urls: 0\nSuccess redirect: 0\nFailed redirect: 0\n"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected body response:\n %s", rr.Body.String())
	}
}

func TestHomeWithoutRedir(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	data := newUrlsStruct()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(data.home)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "This is the home of my website!\n\n"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected body response:\n %s", rr.Body.String())
	}
}

func TestHomeWithRedir(t *testing.T) {
	data := newUrlsStruct()
	url := data.createShortURL("www.test.com")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler2 := http.HandlerFunc(data.home)

	handler2.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "This is the home of my website!\n\nRedirect to:\n" + data.urls[url]
	if rr.Body.String() != expected {
		t.Errorf("Unexpected body response:\n %s", rr.Body.String())
	}
}

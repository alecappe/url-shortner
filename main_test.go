package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func requestHelper(t *testing.T, endpoint string, h http.HandlerFunc) string {
	t.Helper()
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	return rr.Body.String()
}

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
	data := newUrlsStruct()
	res := requestHelper(t, "/shorten/www.google.it", http.HandlerFunc(data.handler))

	if len(res) != 9 {
		t.Errorf("Handler returned unexpected string length")
	}

	if data.urls[res] != "www.google.it" {
		t.Errorf("The saved url in the map is wrong")
	}
}

func TestShowStats(t *testing.T) {
	data := newUrlsStruct()
	res := requestHelper(t, "/stats", http.HandlerFunc(data.showStats))

	expected := "Home called: 0\nShorten called: 0\nStats called: 1\nGenerated urls: 0\nSuccess redirect: 0\nFailed redirect: 0\n"
	if res != expected {
		t.Errorf("Unexpected body response:\n %s", res)
	}
}

func TestHomeWithoutRedir(t *testing.T) {
	data := newUrlsStruct()
	res := requestHelper(t, "/", http.HandlerFunc(data.home))

	expected := "This is the home of my website!\n\n"
	if res != expected {
		t.Errorf("Unexpected body response:\n %s", res)
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

func TestLoadURLHasSucceed(t *testing.T) {
	d := newUrlsStruct()
	data := `{ "urlnum1": "google.it", "urlnum2": "golang.org"}`
	r := strings.NewReader(data)
	if err := d.loadURL(r); err != nil {
		t.Errorf("Load urls failed: %s", err.Error())
	}
}

func TestLoadURLFailed(t *testing.T) {
	d := newUrlsStruct()
	data := "url: google.it"
	r := strings.NewReader(data)
	if err := d.loadURL(r); err == nil {
		t.Errorf("Load urls must be fail")
	}
}

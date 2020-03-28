package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	expectedHeaderNoAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "image/webp,*/*"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", ""},
		{"Cache-Control", "no-cache"},
		{"DNT", "1"},
		{"Pragma", "no-cache"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedHeaderWithAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "image/webp,*/*"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", "Basic bXl1c2VyOm15cGFzc3dvcmQ="},
		{"Cache-Control", "no-cache"},
		{"DNT", "1"},
		{"Pragma", "no-cache"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}
)

func TestDownloadPictureNoAuthorization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

		for _, header := range expectedHeaderNoAuthorization {
			if r.Header.Get(header.name) != header.value {
				t.Fatalf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	err := downloadPicture(ts.URL, "test://referer", "", "Mozilla/5.0 (test)", "", "")
	if err != nil {
		t.Fatalf("downloadPicture returned an error: %v", err)
	}
}

func TestDownloadPictureWithAuthorization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, authorized client")

		for _, header := range expectedHeaderWithAuthorization {
			if r.Header.Get(header.name) != header.value {
				t.Fatalf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	err := downloadPicture(ts.URL, "test://referer", "", "Mozilla/5.0 (test)", "myuser", "mypassword")
	if err != nil {
		t.Fatalf("downloadPicture returned an error: %v", err)
	}
}

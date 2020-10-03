package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testBrowserConfiguration = BrowserConfiguration{
		UserAgent: "Mozilla/5.0 (test)",
		HTML: ElementConfiguration{
			Headers: map[string]string{
				"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Encoding":           "gzip, deflate, br",
				"Accept-Language":           "en-GB,en;q=0.5",
				"Connection":                "keep-alive",
				"DNT":                       "1",
				"Upgrade-Insecure-Requests": "1",
			},
		},
		Picture: ElementConfiguration{
			Headers: map[string]string{
				"Accept":          "image/webp,*/*",
				"Accept-Encoding": "gzip, deflate, br",
				"Accept-Language": "en-GB,en;q=0.5",
				"Cache-Control":   "no-cache",
				"Connection":      "keep-alive",
				"DNT":             "1",
				"Pragma":          "no-cache",
			},
		},
	}

	expectedHTMLHeaderNoAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", ""},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedPictureHeaderNoAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "image/webp,*/*"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", ""},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Pragma", "no-cache"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedPictureHeaderWithAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "image/webp,*/*"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", "Basic bXl1c2VyOm15cGFzc3dvcmQ="},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Pragma", "no-cache"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}
)

func TestDownloadHTMLNoAuthorization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

		for _, header := range expectedHTMLHeaderNoAuthorization {
			if r.Header.Get(header.name) != header.value {
				t.Errorf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0)
	buffer, err := downloadHTML(ts.URL, download)
	if err != nil {
		t.Fatalf("downloadHTML returned an error: %v", err)
	}
	if len(buffer) != 14 {
		t.Fatalf("buffer length should be 14 but returned %d", len(buffer))
	}
}

func TestDownloadPictureNoAuthorization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

		for _, header := range expectedPictureHeaderNoAuthorization {
			if r.Header.Get(header.name) != header.value {
				t.Errorf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0)
	size, err := downloadPicture(ts.URL, "", download)
	if err != nil {
		t.Fatalf("downloadPicture returned an error: %v", err)
	}
	if size != 14 {
		t.Fatalf("size should be 14 but returned %d", size)
	}
}

func TestDownloadPictureWithAuthorization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, authorized client")

		for _, header := range expectedPictureHeaderWithAuthorization {
			if r.Header.Get(header.name) != header.value {
				t.Errorf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	download := NewDownloadConfig(nil, "test://referer", "myuser", "mypassword", "", testBrowserConfiguration, 0, 0)
	size, err := downloadPicture(ts.URL, "", download)
	if err != nil {
		t.Fatalf("downloadPicture returned an error: %v", err)
	}
	if size != 25 {
		t.Fatalf("size should be 25 but returned %d", size)
	}
}

package main

import (
	"fmt"
	"gallery-downloader/config"
	"gallery-downloader/headers"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testBrowserConfiguration = config.Browser{
		Default: config.Group{
			Headers: map[string]string{
				headers.UserAgent:               "Mozilla/5.0 (test)",
				headers.AcceptLanguage:          "en-GB,en;q=0.5",
				headers.DoNotTrack:              "1",
				headers.UpgradeInsecureRequests: "1",
			},
		},
		HTTP: config.Group{
			Headers: map[string]string{
				"Accept-Encoding": "gzip, deflate",
				"Connection":      "keep-alive",
			},
		},
		HTTPS: config.Group{
			Headers: map[string]string{
				"Accept-Encoding": "gzip, deflate, br",
			},
		},
		HTML: config.Group{
			Headers: map[string]string{
				"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			},
		},
		Picture: config.Group{
			Headers: map[string]string{
				"Accept": "image/webp,*/*",
			},
		},
	}

	expectedHTMLHeaderNoAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate"},
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
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", ""},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedPictureHeaderWithAuthorization = []struct {
		name  string
		value string
	}{
		{"Accept", "image/webp,*/*"},
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Authorization", "Basic bXl1c2VyOm15cGFzc3dvcmQ="},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedHTMLHeaderHTTP = []struct {
		name  string
		value string
	}{
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"Connection", "keep-alive"},
		{"DNT", "1"},
		{"Referer", "test://referer"},
		{"User-Agent", "Mozilla/5.0 (test)"},
	}

	expectedHTMLHeaderHTTPS = []struct {
		name  string
		value string
	}{
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-GB,en;q=0.5"},
		{"DNT", "1"},
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

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0, false)
	download.Client = ts.Client()
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

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0, false)
	download.Client = ts.Client()
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

	download := NewDownloadConfig(nil, "test://referer", "myuser", "mypassword", "", testBrowserConfiguration, 0, 0, false)
	download.Client = ts.Client()
	size, err := downloadPicture(ts.URL, "", download)
	if err != nil {
		t.Fatalf("downloadPicture returned an error: %v", err)
	}
	if size != 25 {
		t.Fatalf("size should be 25 but returned %d", size)
	}
}

func TestDownloadHTMLwithHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

		for _, header := range expectedHTMLHeaderHTTP {
			if r.Header.Get(header.name) != header.value {
				t.Errorf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0, false)
	download.Client = ts.Client()
	buffer, err := downloadHTML(ts.URL, download)
	if err != nil {
		t.Fatalf("downloadHTML returned an error: %v", err)
	}
	if len(buffer) != 14 {
		t.Fatalf("buffer length should be 14 but returned %d", len(buffer))
	}
}

func TestDownloadHTMLwithHTTPS(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

		for _, header := range expectedHTMLHeaderHTTPS {
			if r.Header.Get(header.name) != header.value {
				t.Errorf("Incorrect header %s: expected '%s' but found '%s'", header.name, header.value, r.Header.Get(header.name))
			}
		}
	}))
	defer ts.Close()

	download := NewDownloadConfig(nil, "test://referer", "", "", "", testBrowserConfiguration, 0, 0, false)
	download.Client = ts.Client()
	buffer, err := downloadHTML(ts.URL, download)
	if err != nil {
		t.Fatalf("downloadHTML returned an error: %v", err)
	}
	if len(buffer) != 14 {
		t.Fatalf("buffer length should be 14 but returned %d", len(buffer))
	}
}

package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	var err error
	configSource := `{
		"browser": {
			"userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:74.0) Gecko/20100101 Firefox/74.0",
			"html": {
				"headers": {
					"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
					"Accept-Encoding": "gzip, deflate, br",
					"Accept-Language": "en-GB,en;q=0.5",
					"Connection": "keep-alive",
					"DNT": "1",
					"Upgrade-Insecure-Requests": "1"
				}
			},
			"picture": {
				"headers": {
					"Accept": "image/webp,*/*",
					"Accept-Encoding": "gzip, deflate, br",
					"Accept-Language": "en-GB,en;q=0.5",
					"Connection": "keep-alive",
					"DNT": "1",
					"TE": "Trailers"
				}
			}
		}
	}`

	reader := ioutil.NopCloser(bytes.NewReader([]byte(configSource)))
	config, err := loadConfiguration(reader)
	if err != nil {
		t.Fatal(err)
		return
	}

	if config.Browser.UserAgent == "" {
		t.Error("'userAgent' should not be empty")
	}

	if config.Browser.HTML.Headers == nil {
		t.Error("headers configuration not found in section 'html'")
	}
	if len(config.Browser.HTML.Headers) != 6 {
		t.Errorf("'html' section should declare %d headers, but %d found", 6, len(config.Browser.HTML.Headers))
	}
	value, found := config.Browser.HTML.Headers["Upgrade-Insecure-Requests"]
	if !found {
		t.Error("'Upgrade-Insecure-Requests' header not found in 'html' section")
	} else if value != "1" {
		t.Errorf("'Upgrade-Insecure-Requests' header expected value '1' but found %s", value)
	}

	if config.Browser.Picture.Headers == nil {
		t.Error("headers configuration not found in section 'picture'")
	}
	if len(config.Browser.Picture.Headers) != 6 {
		t.Errorf("'picture' section should declare %d headers, but %d found", 6, len(config.Browser.HTML.Headers))
	}
	value, found = config.Browser.Picture.Headers["DNT"]
	if !found {
		t.Error("'DNT' header not found in 'picture' section")
	} else if value != "1" {
		t.Errorf("'DNT' header expected value '1' but found %s", value)
	}
}

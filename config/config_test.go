package config

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var configSource = `{
	"browser": {
		"default": {
			"headers": {
				"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0",
				"Accept-Language": "en-GB,en;q=0.5",
				"DNT": "1",
				"Upgrade-Insecure-Requests": "1"
			}
		},
		"http": {
			"headers": {
				"Accept-Encoding": "gzip, deflate",
				"Connection": "keep-alive"
			}
		},
		"https": {
			"headers": {
				"Accept-Encoding": "gzip, deflate, br"
			}
		},
		"http2": {
			"headers": {
				"Accept-Encoding": "gzip, deflate, br",
				"TE": "Trailers"
			}
		},
		"html": {
			"headers": {
				"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
			}
		},
		"picture": {
			"headers": {
				"Accept": "image/webp,*/*"
			}
		}
	},
	"profiles": [
		{
			"priority": 10,
			"name": "ListItem",
			"generator": "<!--\\s*Generated by\\s*(.*?)\\s*-->",
			"detectGallery": "",
			"detectImage": "<li><img src=\"([^\"]+\\.jpg)\"[^>]*><\/li>",
			"minWait": 0,
			"maxWait": 0,
			"parallel": 5
		},
		{
			"priority": 20,
			"name": "AnchorHREF",
			"generator": "<!--\\s*Generated by\\s*(.*?)\\s*-->",
			"detectGallery": "",
			"detectImage": "<a href=\"([^\"]+\\.jpg)\"[^>]*>.*<\/a>",
			"minWait": 1500,
			"maxWait": 4000,
			"parallel": 1
		}
	]
}`

func TestLoadConfiguration(t *testing.T) {
	var err error

	reader := ioutil.NopCloser(bytes.NewReader([]byte(configSource)))
	cfg, err := loadConfiguration(reader)
	if err != nil {
		t.Fatal(err)
		return
	}

	require.NotEmpty(t, cfg)
	require.NotEmpty(t, cfg.Browser)

	require.NotEmpty(t, cfg.Browser.Default)
	require.NotEmpty(t, cfg.Browser.HTTP)
	require.NotEmpty(t, cfg.Browser.HTTPS)
	require.NotEmpty(t, cfg.Browser.HTTP2)
	require.NotEmpty(t, cfg.Browser.HTML)
	require.NotEmpty(t, cfg.Browser.Picture)

	require.NotEmpty(t, cfg.Browser.Default.Headers)
	require.NotEmpty(t, cfg.Browser.HTTP.Headers)
	require.NotEmpty(t, cfg.Browser.HTTPS.Headers)
	require.NotEmpty(t, cfg.Browser.HTTP2.Headers)
	require.NotEmpty(t, cfg.Browser.HTML.Headers)
	require.NotEmpty(t, cfg.Browser.Picture.Headers)

	// Check a few random values
	assert.NotEmpty(t, cfg.Browser.Default.Headers["User-Agent"])
	assert.NotEmpty(t, cfg.Browser.HTTP.Headers["Accept-Encoding"])
	assert.NotEmpty(t, cfg.Browser.HTML.Headers["Accept"])

	assert.NotEmpty(t, cfg.Profiles)
	for _, profile := range cfg.Profiles {
		if profile.DetectGallery != "" {
			_, err := regexp.Compile(profile.DetectGallery)
			assert.NoErrorf(t, err, "cannot parse `%s`", profile.DetectGallery)
		}
		if profile.DetectImage != "" {
			_, err := regexp.Compile(profile.DetectImage)
			assert.NoErrorf(t, err, "cannot parse `%s`", profile.DetectImage)
		}
		if profile.Generator != "" {
			_, err := regexp.Compile(profile.Generator)
			assert.NoErrorf(t, err, "cannot parse `%s`", profile.Generator)
		}
	}
}

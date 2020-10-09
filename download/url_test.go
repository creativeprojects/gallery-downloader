package download

import (
	"net/url"
	"testing"
)

var (
	testURLs = []struct {
		first    string
		second   string
		expected string
	}{
		{
			"http://localhost",
			"dir/file",
			"http://localhost/dir/file",
		},
		{
			"http://localhost/",
			"file",
			"http://localhost/file",
		},
		{
			"http://localhost/base",
			"file",
			"http://localhost/file",
		},
		{
			"http://localhost/base/",
			"file",
			"http://localhost/base/file",
		},
		{
			"http://localhost/base",
			"/file",
			"http://localhost/file",
		},
		{
			"http://localhost/base/",
			"/dir/file",
			"http://localhost/dir/file",
		},
		{
			"http://localhost/base/index",
			"/dir/file",
			"http://localhost/dir/file",
		},
		{
			"http://localhost/base/index",
			"dir/file",
			"http://localhost/base/dir/file",
		},
	}
)

func TestJoinURLs(t *testing.T) {
	for index, test := range testURLs {
		first, _ := url.Parse(test.first)
		second, _ := url.Parse(test.second)
		result := joinURL(first, second)
		if result.String() != test.expected {
			t.Errorf("Test %d: Expected '%s' but found '%s' (joining '%s' with '%s')", index+1, test.expected, result.String(), first.String(), second.String())
		}
	}
}

package scan

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpMatcherFind(t *testing.T) {
	pattern := regexp.MustCompile("<!--\\s*Generated by\\s*(.*?)\\s*-->")
	var matcher Matcher = NewRegexpMatcher(pattern)
	assert.Equal(t, "", matcher.Find([]byte("blahblahblah")))
	assert.Equal(t, "WOWSlider.com v5.5", matcher.Find([]byte(testHTMLtype1)))
	assert.Equal(t, "WOWSlider.com v5.6", matcher.Find([]byte(testHTMLtype2)))
}

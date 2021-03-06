package scan

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpMatcherFind(t *testing.T) {
	pattern := regexp.MustCompile(`<!--\s*Generated by\s*(.*?)\s*-->`)
	var matcher Matcher = NewRegexpMatcher(pattern)
	matcher.Source([]byte("blahblahblah"))
	assert.Equal(t, "", matcher.Find())

	matcher.Source(getTestData(t, "anchor_href"))
	assert.Equal(t, "WOWSlider.com v5.5", matcher.Find())

	matcher.Source(getTestData(t, "list_item"))
	assert.Equal(t, "WOWSlider.com v5.6", matcher.Find())
}

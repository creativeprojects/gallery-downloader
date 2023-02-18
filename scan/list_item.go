package scan

import (
	"bytes"
	"regexp"

	"golang.org/x/net/html"
)

// LegacyListItemGallery scans pictures like <li><img src="data1/images/picture002.jpg" alt="picture-002" title="picture-002" id="wows1_1"/></li>
type LegacyListItemGallery struct {
	source []byte
	node   *html.Node
}

// NewLegacyListItemGallery creates a new gallery
func NewLegacyListItemGallery(source []byte) (Gal, error) {
	node, err := html.Parse(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}
	return &LegacyListItemGallery{
		source: source,
		node:   node,
	}, nil
}

// HasDetection returns true when the current type of gallery can be detected
func (g *LegacyListItemGallery) HasDetection() bool {
	return false
}

// Match returns true if this profile *can* be a match for the current file.
// if there's no gallery detection, it returns true to try to find images
func (g *LegacyListItemGallery) Match() bool {
	return true
}

// GeneratedBy returns the name of the gallery generator (if available).
// It returns an empty string if not available
func (g *LegacyListItemGallery) GeneratedBy() string {
	pattern := regexp.MustCompile(`<!--\s*Generated by\s*(.*?)\s*-->`)
	match := pattern.FindSubmatch(g.source)
	if match == nil || len(match) != 2 {
		return ""
	}
	return string(match[1])
}

// Find returns a list of images found in this gallery
func (g *LegacyListItemGallery) Find() []string {
	pictures := make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if picture, found := getPictureAttribute(c, "img", "src"); found {
					pictures = append(pictures, picture)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(g.node)
	return pictures
}

// Verify interfaces
var (
	_ Gal        = &LegacyListItemGallery{}
	_ GalFactory = NewLegacyListItemGallery
)
package scan

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// GalleryScannerFunc is a gallery scanner function returning picture URLs
type GalleryScannerFunc func(source io.Reader) ([]string, error)

// Gallery type
const (
	AutoDetect = "AutoDetect"
	AnchorHREF = "AnchorHREF"
	ListItem   = "ListItem"
)

var (
	// AvailableGalleryScanners lists the available gallery scanners
	AvailableGalleryScanners = [...]string{AutoDetect, AnchorHREF, ListItem}
	// GalleryScanners maps the gallery scanner functions
	GalleryScanners map[string][]GalleryScannerFunc
)

func init() {
	GalleryScanners = map[string][]GalleryScannerFunc{
		AutoDetect: {
			GalleryAnchorHREF,
			GalleryListItem,
		},
		AnchorHREF: {GalleryAnchorHREF},
		ListItem:   {GalleryListItem},
	}
}

// GalleryAnchorHREF scans pictures like <a href="picture2.jpg" title="picture2">picture 2</a>
func GalleryAnchorHREF(source io.Reader) ([]string, error) {
	doc, err := html.Parse(source)
	if err != nil {
		return nil, err
	}
	pictures := make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if picture, found := getPictureAttribute(n, "a", "href"); found {
			pictures = append(pictures, picture)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return pictures, nil
}

// GalleryListItem scans pictures like <li><img src="data1/images/picture002.jpg" alt="picture-002" title="picture-002" id="wows1_1"/></li>
func GalleryListItem(source io.Reader) ([]string, error) {
	doc, err := html.Parse(source)
	if err != nil {
		return nil, err
	}
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
	f(doc)
	return pictures, nil
}

func getPictureAttribute(n *html.Node, element, attribute string) (string, bool) {
	if n.Type == html.ElementNode && n.Data == element {
		for _, a := range n.Attr {
			if a.Key == attribute {
				value := strings.ToLower(a.Val)
				if strings.HasSuffix(value, "jpg") || strings.HasSuffix(value, "jpeg") {
					return a.Val, true
				}
				break
			}
		}
	}
	return "", false
}
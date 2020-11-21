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
	AutoDetect     = "AutoDetect"
	AnchorHREF     = "AnchorHREF"
	ListItem       = "ListItem"
	ConfigProfiles = "ConfigProfiles"
)

var (
	// AvailableGalleryScanners lists the available gallery scanners
	AvailableGalleryScanners = [...]string{AutoDetect, ConfigProfiles, AnchorHREF, ListItem}
	// GalleryScanners maps the gallery scanner constructors
	GalleryScanners map[string][]GalFactory
)

func init() {
	GalleryScanners = map[string][]GalFactory{
		AnchorHREF: {NewLegacyAnchorGallery},
		ListItem:   {NewLegacyListItemGallery},
	}
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

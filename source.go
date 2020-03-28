package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func loadGalleryAnchorHREF(source io.Reader) ([]string, error) {
	doc, err := html.Parse(source)
	if err != nil {
		return nil, err
	}
	pictures := make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					value := strings.ToLower(a.Val)
					if strings.HasSuffix(value, "jpg") || strings.HasSuffix(value, "jpeg") {
						pictures = append(pictures, a.Val)
					}
					break
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

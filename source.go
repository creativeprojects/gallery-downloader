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

func loadGalleryListItem(source io.Reader) ([]string, error) {
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

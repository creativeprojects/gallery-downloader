package scan

import (
	"bytes"
	"log"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type SelectorMatcher struct {
	sel       cascadia.Sel
	attribute string
	source    *html.Node
}

func NewSelectorMatcher(sel cascadia.Sel, attribute string) *SelectorMatcher {
	if sel == nil {
		// might as well panic right now, no need to go much further
		panic("invalid nil css selector pattern")
	}
	return &SelectorMatcher{
		sel:       sel,
		attribute: attribute,
	}
}

func (m *SelectorMatcher) Source(source []byte) error {
	var err error

	m.source, err = html.Parse(bytes.NewReader(source))
	if err != nil {
		return err
	}
	return nil
}

func (m *SelectorMatcher) Find() string {
	node := cascadia.Query(m.source, m.sel)
	if node == nil {
		return ""
	}
	buffer := &strings.Builder{}
	html.Render(buffer, node)
	log.Printf("Find(): %q", buffer.String())
	return buffer.String()
}

func (m *SelectorMatcher) FindAll() []string {
	nodes := cascadia.QueryAll(m.source, m.sel)
	if nodes == nil {
		return nil
	}
	images := make([]string, len(nodes))
	for i, node := range nodes {
		images[i] = getAttribute(node, m.attribute)
	}
	log.Printf("FindAll(): %v", images)
	return images
}

func getAttribute(n *html.Node, attribute string) string {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == attribute {
				return a.Val
			}
		}
	}
	return ""
}

// Verify interface
var _ Matcher = &SelectorMatcher{}

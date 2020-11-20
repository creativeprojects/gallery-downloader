package scan

import "regexp"

type RegexpMatcher struct {
	pattern *regexp.Regexp
	source  []byte
}

func NewRegexpMatcher(pattern *regexp.Regexp) *RegexpMatcher {
	if pattern == nil {
		// might as well panic right now, no need to go much further
		panic("invalid nil regexp pattern")
	}
	return &RegexpMatcher{
		pattern: pattern,
	}
}

func (m *RegexpMatcher) Source(source []byte) error {
	m.source = source
	return nil
}

func (m *RegexpMatcher) Find() string {
	found := m.pattern.FindSubmatch(m.source)
	if found == nil {
		return ""
	}
	if len(found) >= 2 {
		// always return the first submatch ...
		return string(found[1])
	}
	// ... or the whole expression if there was no catching parenthesis
	return string(found[0])
}

func (m *RegexpMatcher) FindAll() []string {
	all := m.pattern.FindAllSubmatch(m.source, -1)
	if all == nil {
		return nil
	}
	found := make([]string, len(all))
	for i, match := range all {
		if len(match) >= 2 {
			// always return the first submatch ...
			found[i] = string(match[1])
			continue
		}
		// ... or the whole expression if there was no catching parenthesis
		found[i] = string(match[0])
	}
	return found
}

// Verify interface
var _ Matcher = &RegexpMatcher{}

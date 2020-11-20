package scan

// Gallery profile to detect images
type Gallery struct {
	cfg    Config
	source []byte
}

// NewGallery creates a new profile of gallery
func NewGallery(cfg Config, source []byte) *Gallery {
	return &Gallery{
		cfg:    cfg,
		source: source,
	}
}

// HasDetection returns true when the current type of gallery can be detected
func (g *Gallery) HasDetection() bool {
	return g.cfg.DetectGallery != nil
}

// Match returns true if this profile *can* be a match for the current file.
// if there's no gallery detection, it returns true to try to find images
func (g *Gallery) Match() bool {
	return g.cfg.DetectGallery == nil || g.cfg.DetectGallery.FindIndex(g.source) != nil
}

// GeneratedBy returns the name of the gallery generator (if available).
// It returns an empty string if not available
func (g *Gallery) GeneratedBy() string {
	if g.cfg.DetectGenerator == nil {
		return ""
	}
	match := g.cfg.DetectGenerator.FindSubmatch(g.source)
	if match == nil || len(match) != 2 {
		return ""
	}
	return string(match[1])
}

// Found returns a list of images found in this gallery
func (g *Gallery) Found() []string {
	all := g.cfg.DetectImage.FindAllSubmatch(g.source, -1)
	if all == nil {
		return nil
	}
	found := make([]string, len(all))
	for i, match := range all {
		if len(match) != 2 {
			continue
		}
		found[i] = string(match[1])
	}
	return found
}

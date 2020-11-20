package scan

// Gallery profile to detect images
type Gallery struct {
	cfg Config
	// source []byte
}

// NewGallery creates a new profile of gallery
func NewGallery(cfg Config, source []byte) *Gallery {
	if cfg.DetectGenerator != nil {
		cfg.DetectGenerator.Source(source)
	}
	if cfg.DetectGallery != nil {
		cfg.DetectGallery.Source(source)
	}
	if cfg.DetectImage != nil {
		cfg.DetectImage.Source(source)
	}
	return &Gallery{
		cfg: cfg,
		// source: source,
	}
}

// HasDetection returns true when the current type of gallery can be detected
func (g *Gallery) HasDetection() bool {
	return g.cfg.DetectGallery != nil
}

// Match returns true if this profile *can* be a match for the current file.
// if there's no gallery detection, it returns true to try to find images
func (g *Gallery) Match() bool {
	return g.cfg.DetectGallery == nil || g.cfg.DetectGallery.Find() != ""
}

// GeneratedBy returns the name of the gallery generator (if available).
// It returns an empty string if not available
func (g *Gallery) GeneratedBy() string {
	if g.cfg.DetectGenerator == nil {
		return ""
	}
	return g.cfg.DetectGenerator.Find()
}

// Found returns a list of images found in this gallery
func (g *Gallery) Found() []string {
	return g.cfg.DetectImage.FindAll()
}

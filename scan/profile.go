package scan

type Gallery struct {
	cfg    Config
	source []byte
}

func NewGallery(cfg Config, source []byte) *Gallery {
	return &Gallery{
		cfg:    cfg,
		source: source,
	}
}

func (g *Gallery) HasDetection() bool {
	return g.cfg.DetectGallery != nil
}

func (g *Gallery) IsDetected() bool {
	if g.cfg.DetectGallery == nil {
		panic("no detection available, please check HasDetection() first")
	}
	return g.cfg.DetectGallery.FindIndex(g.source) != nil
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

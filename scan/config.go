package scan

import "regexp"

// Config contains the gallery profile configuration
type Config struct {
	Name            string
	DetectGenerator *regexp.Regexp
	DetectGallery   *regexp.Regexp
	DetectImage     *regexp.Regexp
}

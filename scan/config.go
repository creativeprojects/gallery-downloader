package scan

import "regexp"

type Config struct {
	Name            string
	DetectGenerator *regexp.Regexp
	DetectGallery   *regexp.Regexp
	DetectImage     *regexp.Regexp
}

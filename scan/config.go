package scan

// Config contains the gallery profile configuration
type Config struct {
	Name            string
	DetectGenerator Matcher
	DetectGallery   Matcher
	DetectImage     Matcher
}

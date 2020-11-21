package scan

// GalFactory is a concrete constructor for an object providing a Gal interface
type GalFactory func(source []byte) (Gal, error)

// Gal is a gallery interface
type Gal interface {
	HasDetection() bool
	Match() bool
	GeneratedBy() string
	Find() []string
}

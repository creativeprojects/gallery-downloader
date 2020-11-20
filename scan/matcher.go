package scan

type Matcher interface {
	Source(source []byte) error
	Find() string
	FindAll() []string
}

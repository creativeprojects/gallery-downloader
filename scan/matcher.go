package scan

type Matcher interface {
	Find(source []byte) string
	FindAll(source []byte) []string
}

package parsers

type Parser interface {
	Parse(text string) (map[string]string, error)
}

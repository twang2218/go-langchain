package parsers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexParser(t *testing.T) {
	var p Parser = &RegexParser{
		Regex: regexp.MustCompile(`(?i)hello (.*)`),
		Keys:  []string{"name"},
	}

	result, err := p.Parse("Hello World")
	assert.NoError(t, err)
	assert.Equal(t, "World", result["name"])
}

package parsers

import (
	"fmt"
	"regexp"
)

type RegexParser struct {
	Regex *regexp.Regexp
	Keys  []string
}

func (p *RegexParser) Parse(output string) (map[string]string, error) {
	match := p.Regex.FindStringSubmatch(output)
	if len(match) != len(p.Keys)+1 {
		return nil, fmt.Errorf("could not parse output: %s", output)
	}

	result := make(map[string]string)
	for i, key := range p.Keys {
		result[key] = match[i+1]
	}

	return result, nil
}

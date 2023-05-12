package prompts

import (
	"regexp"
	"strings"
)

type PromptTemplate struct {
	Template  string
	Variables []string
}

func NewPromptTemplateByTemplate(template string) PromptTemplate {
	reVariables := regexp.MustCompile(`\{([^\{\}]+)\}`)
	var variables []string
	for _, match := range reVariables.FindAllStringSubmatch(template, -1) {
		v := strings.TrimSpace(match[1])
		variables = append(variables, v)
	}
	return PromptTemplate{
		Template:  template,
		Variables: variables,
	}
}

func (t PromptTemplate) Format(values map[string]string) string {
	prompt := t.Template
	for _, variable := range t.Variables {
		if v, ok := values[variable]; ok {
			prompt = strings.ReplaceAll(prompt, "{"+variable+"}", v)
		}
	}
	return prompt
}

func (t PromptTemplate) GetVariables() []string {
	return t.Variables
}

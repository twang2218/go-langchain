package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptTemplate(t *testing.T) {
	var template Template = NewPromptTemplateByTemplate("Hello {name}!")
	vars := template.GetVariables()
	assert.Equal(t, 1, len(vars))
	assert.Equal(t, "name", vars[0])
	prompt := template.Format(map[string]string{
		"name": "World",
	})
	assert.Equal(t, "Hello World!", prompt)
	prompt = template.Format(map[string]string{
		"name": "China",
	})
	assert.Equal(t, "Hello China!", prompt)
}

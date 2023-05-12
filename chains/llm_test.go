package chains

import (
	"context"
	"langchain/llms"
	"langchain/prompts"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLLMBasic(t *testing.T) {
	var llm llms.LLM = &llms.OpenAI{}
	// logrus.SetLevel(logrus.TraceLevel)

	//	Basic template
	var c Chain = NewLLMChain(llm, nil)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*LLMChain).LLM)
	assert.NotNil(t, c.(*LLMChain).Prompt)

	query := "one plus one equals two, is it true? please answer only using number, 1 for yes, 0 for no"

	ctx := context.Background()
	resp, err := c.RunText(ctx, query)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "1", resp)

	resp, err = c.Run(ctx, map[string]string{
		"prompt": query,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "1", resp)

	resps, err := c.Apply(ctx, []map[string]string{
		{
			"prompt": query,
		},
		{
			"prompt": query,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resps)
	assert.Len(t, resps, 2, "should be 2 responses: %v", resps)
	t.Logf("resps: %v", resps)
	if len(resps) == 2 {
		assert.Regexp(t, "1", resps[0])
		assert.Regexp(t, "1", resps[1])
	}
}

func TestLLMCustom(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	//	Basic template

	//	Template with variables
	template := "one plus one equals two, is it true? please answer only using single {quantity}, {positive} for yes, {negative} for no"
	var c Chain = NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(template))
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*LLMChain).LLM)
	assert.NotNil(t, c.(*LLMChain).Prompt)
	assert.Len(t, c.(*LLMChain).Prompt.GetVariables(), 3)

	resp, err := c.RunText(ctx, template)
	assert.Error(t, err)
	assert.Empty(t, resp)

	resp, err = c.Run(ctx, map[string]string{
		"quantity": "letter",
		"positive": "y",
		"negative": "n",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "y", strings.ToLower(resp))

	resp, err = c.Run(ctx, map[string]string{
		"quantity": "number",
		"positive": "1",
		"negative": "0",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "1", resp)

	resps, err := c.Apply(ctx, []map[string]string{
		{
			"quantity": "number",
			"positive": "1",
			"negative": "0",
		},
		{
			"quantity": "letter",
			"positive": "y",
			"negative": "n",
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resps)
	assert.Len(t, resps, 2, "should be 2 responses: %v", resps)
	if len(resps) == 2 {
		assert.Regexp(t, "1", resps[0])
		assert.Regexp(t, "y", strings.ToLower(resps[1]))
	}
}

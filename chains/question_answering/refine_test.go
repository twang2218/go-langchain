package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"langchain/prompts"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefineBasic(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Basic
	var c chains.Chain = NewRefineDocumentsChain(llm, nil, nil)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain.(*chains.LLMChain).Prompt)

	values := map[string]string{}
	docs := []string{
		"Alice is older than Bob",
		"Bob is older than Charlie",
		"Charlie likes to play chess",
		"Bob drinks coffee a lot",
	}
	err := PutDocsToInputs(docs, values, KeyDocuments)
	assert.NoError(t, err)
	assert.NotEmpty(t, values[KeyDocuments])
	values[KeyQuestion] = "Is Alice older than Bob?"

	resp, err := c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "(yes|alice is older than bob)", strings.ToLower(resp))
}

func TestRefineCustom(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Custom
	initialChain := chains.NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(`
	用户提出的问题：
	------
	{`+KeyQuestion+`}
	------
	现在有一些信息：
	------
	{`+KeyContext+`}
	------
	答案:`))
	refineChain := chains.NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(`
	对于用户提出的问题：
	------
	{`+KeyQuestion+`}
	------
	现有的回答是：
	------
	{`+KeyExistingAnswer+`}
	------
	现在有一些新的信息：
	------
	{`+KeyContext+`}
	------
	请问现在这个问题的答案是什么？如果新的信息不足以回答这个问题或与问题无关，请直接回复现有的回答。
	最终答案为:`))
	var c chains.Chain = NewRefineDocumentsChain(llm, initialChain, refineChain)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*RefineDocumentsChain).InitialChain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*RefineDocumentsChain).RefineChain.(*chains.LLMChain).Prompt)

	values := map[string]string{}
	docs := []string{
		"Alice is older than Bob",
		"Bob is older than Charlie",
		"Charlie likes to play chess",
		"Bob drinks coffee a lot",
	}
	err := PutDocsToInputs(docs, values, KeyDocuments)
	assert.NoError(t, err)
	assert.NotEmpty(t, values[KeyDocuments])
	values[KeyQuestion] = "Is Alice older than Bob?"

	resp, err := c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "(yes|是的|正确)", strings.ToLower(resp))
}

package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"langchain/prompts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStuffDocumentsChainBasic(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Basic template
	var c chains.Chain = NewStuffDocumentsChain(llm, nil)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain.(*chains.LLMChain).Prompt)

	values := map[string]string{}
	docs := []string{
		"Alice is older than Bob",
		"Bob is older than Charlie",
	}
	err := PutDocsToInputs(docs, values, KeyDocuments)
	assert.NoError(t, err)
	assert.NotEmpty(t, values[KeyDocuments])
	values[KeyQuestion] = "Who is the oldest?"

	resp, err := c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "Alice", resp)

	resp, err = c.RunText(ctx, "Who is the oldest?")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.NotRegexp(t, "Alice", resp)
}

func TestStuffDocumentsChainCustom(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Custom template
	template := `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

	{` + KeyContext + `}
	
	Question: {` + KeyQuestion + `}
	Helpful Answer in Chinese:`
	var c chains.Chain = NewStuffDocumentsChain(llm, prompts.NewPromptTemplateByTemplate(template))
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*StuffDocumentsChain).Chain.(*chains.LLMChain).Prompt)

	values := map[string]string{}
	docs := []string{
		"张三喜欢吃西瓜",
		"李四喜欢吃木瓜",
		"张三喜欢打篮球",
		"王五喜欢冲浪",
	}
	err := PutDocsToInputs(docs, values, KeyDocuments)
	assert.NoError(t, err)
	assert.NotEmpty(t, values[KeyDocuments])
	values[KeyQuestion] = "谁喜欢吃西瓜？"

	resp, err := c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "张三", resp)

	values[KeyQuestion] = "有喜欢打篮球的人吗？"
	resp, err = c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "有", resp)

	values[KeyQuestion] = "喜欢运动的有哪些人？"
	resp, err = c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "张三", resp)
	assert.Regexp(t, "王五", resp)
}

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

func TestMapReduceBasic(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Basic
	var c chains.Chain = NewMapReduceDocumentsChain(llm, nil, nil)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).DocumentSeparator)

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

func TestMapReduceCustomZhEn(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Custom templates
	mapPrompt := `使用下列长文档的段落来观察是否有任何文本与问题相关，并逐字返回相关文本内容。
	{` + KeyContext + `}
	问题：{` + KeyQuestion + `}
	相关的文本内容（如果存在的话）；否则只回答'无'即可：`

	reducePrompt := `给定一份长文档的提取部分和一个问题，创建一个最终答案。如果你不知道答案，就说你不知道。不要试图编造答案。
	============
	{` + KeyContext + `}
	============
	问题：{` + KeyQuestion + `}
	答案：`

	var c chains.Chain = NewMapReduceDocumentsChain(
		llm,
		prompts.NewPromptTemplateByTemplate(mapPrompt),
		prompts.NewPromptTemplateByTemplate(reducePrompt))
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).DocumentSeparator)

	values := map[string]string{}
	docs := []string{
		"Bob learnt how to drive car in town, Alice is older than Bob, and Alice run faster than Bob",
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
	assert.Regexp(t, "(yes|是的)", strings.ToLower(resp))
}

func TestMapReduceCustomZhZh(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Custom templates
	mapPrompt := `使用下列长文档的段落来观察是否有任何文本与问题相关，并逐字返回相关文本内容。
	{` + KeyContext + `}
	问题：{` + KeyQuestion + `}
	相关的文本内容（如果存在的话）；否则只回答'无'即可：`

	reducePrompt := `给定一份长文档的提取部分和一个问题，创建一个最终答案。如果你不知道答案，就说你不知道。不要试图编造答案。
	============
	{` + KeyContext + `}
	============
	问题：{` + KeyQuestion + `}
	答案：`

	var c chains.Chain = NewMapReduceDocumentsChain(
		llm,
		prompts.NewPromptTemplateByTemplate(mapPrompt),
		prompts.NewPromptTemplateByTemplate(reducePrompt))
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).MapChain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).Chain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapReduceDocumentsChain).ReduceChain.(*StuffDocumentsChain).DocumentSeparator)

	values := map[string]string{}
	docs := []string{
		"张三来自南京，张三喜欢吃西瓜，西瓜是一种有很多水分的水果",
		"李四喜欢吃木瓜，因为他觉得木瓜有很多维生素，维生素对健康有好处",
		"篮球场上人很多，张三在那里，因为张三喜欢打篮球",
		"平常你是找不到王五的，王五喜欢踢足球，他经常去操场踢足球",
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

	values[KeyQuestion] = "请列出资料中所有喜欢运动的人。"
	resp, err = c.Run(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Regexp(t, "张三", resp)
	assert.Regexp(t, "王五", resp)
}

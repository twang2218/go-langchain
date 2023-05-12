package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapRerankDocumentsChainDefault(t *testing.T) {
	ctx := context.Background()
	var llm llms.LLM = &llms.OpenAI{}

	// logrus.SetLevel(logrus.TraceLevel)

	//	Basic
	var c chains.Chain = NewMapRerankDocumentsChainDefault(llm)
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).Chain)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).Chain.(*chains.LLMChain).LLM)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).Chain.(*chains.LLMChain).Prompt)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).Parser)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).RankKey)
	assert.NotNil(t, c.(*MapRerankDocumentsChain).AnswerKey)

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

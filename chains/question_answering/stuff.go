package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"langchain/prompts"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// Stuff Documents Chain Prompt Template
	DefaultStuffPromptTemplate = `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

	{` + KeyContext + `}
	
	Question: {` + KeyQuestion + `}
	Helpful Answer:`
	DefaultStuffChatPromptTemplate = `Use the following pieces of context to answer the users question. 
	If you don't know the answer, just say that you don't know, don't try to make up an answer.
	----------------
	{` + KeyContext + `}`
)

type StuffDocumentsChain struct {
	Chain             chains.Chain
	DocumentSeparator string
}

func NewStuffDocumentsChain(llm llms.LLM, prompt prompts.Template) *StuffDocumentsChain {
	if prompt == nil {
		prompt = prompts.NewPromptTemplateByTemplate(DefaultStuffPromptTemplate)
	}

	return &StuffDocumentsChain{
		Chain:             chains.NewLLMChain(llm, prompt),
		DocumentSeparator: DefaultDocumentSeparator,
	}
}

func (c *StuffDocumentsChain) Run(ctx context.Context, inputs map[string]string) (string, error) {
	logrus.Tracef("StuffDocumentsChain.Run(): %v", inputs)
	docs, err := GetDocsFromInputs(inputs, KeyDocuments)
	if err != nil {
		return "", err
	}

	// make a copy of inputs without documents
	values := make(map[string]string)
	for k, v := range inputs {
		if k == KeyDocuments {
			continue
		}
		values[k] = v
	}

	return c.CombineDocs(ctx, docs, values)
}

func (c *StuffDocumentsChain) CombineDocs(ctx context.Context, docs []string, inputs map[string]string) (string, error) {
	//	将文档列表合并成一个文档
	var combined_docs string
	if len(docs) == 0 {
		combined_docs = ""
	} else {
		if c.DocumentSeparator == "" {
			c.DocumentSeparator = DefaultDocumentSeparator
		}
		combined_docs = strings.Join(docs, c.DocumentSeparator)
	}
	//	一次性调用LLM获取回答
	//  TODO: 将 KeyContext 改为可以设置的变量，这样可以使用 KeySummaries 之类的其它键值。
	inputs[KeyContext] = combined_docs
	return c.Chain.Run(ctx, inputs)
}

func (c *StuffDocumentsChain) RunText(ctx context.Context, input string) (string, error) {
	inputs := make(map[string]string)
	inputs[KeyDocuments] = "[]"
	inputs[KeyQuestion] = input
	return c.Run(ctx, inputs)
}

func (c *StuffDocumentsChain) Apply(ctx context.Context, inputs []map[string]string) ([]string, error) {
	return chains.ChainApplyAsync(c, ctx, inputs)
}

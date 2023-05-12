package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"langchain/prompts"

	"github.com/sirupsen/logrus"
)

const (
	DefaultRefineInitialPromptTemplate = `Context information is below. 
---------------------
{` + KeyContext + `}
---------------------
Given the context information and not prior knowledge, answer the question: {` + KeyQuestion + `}
Helpful Answer:
`
	DefaultRefinePromptTemplate = `The original question is as follows: {` + KeyQuestion + `}
We have provided an existing answer: {` + KeyExistingAnswer + `}
We have the opportunity to refine the existing answer (only if needed) with some more context below.
------------
{` + KeyContext + `}
------------
Given the new context, refine the original answer to better answer the question. If the context isn't useful, return the original answer.
Answer:`
)

type RefineDocumentsChain struct {
	InitialChain chains.Chain
	RefineChain  chains.Chain
}

func NewRefineDocumentsChain(llm llms.LLM, initialChain, refineChain chains.Chain) *RefineDocumentsChain {
	if initialChain == nil {
		initialChain = chains.NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(DefaultRefineInitialPromptTemplate))
	}
	if refineChain == nil {
		refineChain = chains.NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(DefaultRefinePromptTemplate))
	}

	return &RefineDocumentsChain{
		InitialChain: initialChain,
		RefineChain:  refineChain,
	}
}

func (c *RefineDocumentsChain) Run(ctx context.Context, inputs map[string]string) (string, error) {
	logrus.Tracef("RefineDocumentsChain.Run(): %v", inputs)
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

func (c *RefineDocumentsChain) CombineDocs(ctx context.Context, docs []string, inputs map[string]string) (string, error) {
	//	在没有前一次查询结果的情况下，调用 initialChain 来生成一个初始的答案
	if len(docs) == 0 {
		inputs[KeyContext] = ""
		return c.InitialChain.Run(ctx, inputs)
	}

	//	first doc
	inputs[KeyContext] = docs[0]
	existing_answer, err := c.InitialChain.Run(ctx, inputs)
	if err != nil {
		return "", err
	}

	//	remaining docs
	for _, doc := range docs[1:] {
		inputs[KeyContext] = doc
		inputs[KeyExistingAnswer] = existing_answer
		existing_answer, err = c.RefineChain.Run(ctx, inputs)
		if err != nil {
			return "", err
		}
	}

	return existing_answer, nil
}

func (c *RefineDocumentsChain) RunText(ctx context.Context, input string) (string, error) {
	inputs := make(map[string]string)
	inputs[KeyDocuments] = "[]"
	return c.Run(ctx, inputs)
}

func (c *RefineDocumentsChain) Apply(ctx context.Context, inputs []map[string]string) ([]string, error) {
	return chains.ChainApplyAsync(c, ctx, inputs)
}

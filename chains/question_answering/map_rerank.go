package question_answering

import (
	"context"
	"fmt"
	"langchain/chains"
	"langchain/llms"
	"langchain/parsers"
	"langchain/prompts"
	"regexp"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMapRerankPromptTemplate = `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

In addition to giving an answer, also return a score of how fully it answered the user's question. This should be in the following format:

Question: [question here]
Helpful Answer: [answer here]
Score: [score between 0 and 100]

How to determine the score:
- Higher is a better answer
- Better responds fully to the asked question, with sufficient level of detail
- If you do not know the answer based on the context, that should be a score of 0
- Don't be overconfident!

Example #1

Context:
---------
Apples are red
---------
Question: what color are apples?
Helpful Answer: red
Score: 100

Example #2

Context:
---------
it was night and the witness forgot his glasses. he was not sure if it was a sports car or an suv
---------
Question: what type was the car?
Helpful Answer: a sports car or an suv
Score: 60

Example #3

Context:
---------
Pears are either red or orange
---------
Question: what color are apples?
Helpful Answer: This document does not answer the question
Score: 0

Begin!

Context:
---------
{` + KeyContext + `}
---------
Question: {` + KeyQuestion + `}
Helpful Answer:`
)

type MapRerankDocumentsChain struct {
	Chain     chains.Chain
	Parser    parsers.Parser
	RankKey   string
	AnswerKey string
}

func NewMapRerankDocumentsChainDefault(llm llms.LLM) *MapRerankDocumentsChain {
	return NewMapRerankDocumentsChain(llm, "", "", "", "")
}

func NewMapRerankDocumentsChain(llm llms.LLM, prompt, regex, rank_key, answer_key string) *MapRerankDocumentsChain {
	if prompt == "" {
		prompt = DefaultMapRerankPromptTemplate
	}
	if regex == "" {
		regex = `(?s)(.*?)\nScore: (\d+)`
	}
	if rank_key == "" {
		rank_key = KeyRank
	}
	if answer_key == "" {
		answer_key = KeyAnswer
	}

	parser := &parsers.RegexParser{
		Regex: regexp.MustCompile(regex),
		Keys:  []string{answer_key, rank_key}} // order matters

	return &MapRerankDocumentsChain{
		Chain:     chains.NewLLMChain(llm, prompts.NewPromptTemplateByTemplate(prompt)),
		Parser:    parser,
		RankKey:   rank_key,
		AnswerKey: answer_key,
	}
}

func (c *MapRerankDocumentsChain) Run(ctx context.Context, inputs map[string]string) (string, error) {
	logrus.Tracef("MapRerankDocumentsChain.Run(): %v", inputs)
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

func (c *MapRerankDocumentsChain) CombineDocs(ctx context.Context, docs []string, inputs map[string]string) (string, error) {
	if len(docs) == 0 {
		return "", fmt.Errorf("no documents to combine")
	}

	//  Map each document to the question
	docs_inputs := make([]map[string]string, 0, len(docs))
	for _, doc := range docs {
		//	make a copy of the inputs
		doc_inputs := make(map[string]string)
		for k, v := range inputs {
			doc_inputs[k] = v
		}
		doc_inputs[KeyContext] = doc
		docs_inputs = append(docs_inputs, doc_inputs)
	}

	results, err := c.Chain.Apply(ctx, docs_inputs)
	if err != nil {
		return "", err
	}

	parsed_results := []map[string]string{}
	for _, result := range results {
		if c.Parser == nil {
			return "", fmt.Errorf("no parser for MapRerankDocumentsChain, we need a parser for score to rank")
		}
		parsed_result, err := c.Parser.Parse(result)
		if err != nil {
			return "", err
		}
		parsed_results = append(parsed_results, parsed_result)
	}

	//	sort by score
	sort.Slice(parsed_results, func(i, j int) bool {
		//  convert score to float
		score_i, _ := strconv.Atoi(parsed_results[i][c.RankKey])
		score_j, _ := strconv.Atoi(parsed_results[j][c.RankKey])
		return score_i > score_j
	})

	//	get top one as answer
	return parsed_results[0][c.AnswerKey], nil
}

func (c *MapRerankDocumentsChain) RunText(ctx context.Context, text string) (string, error) {
	inputs := make(map[string]string)
	inputs[KeyDocuments] = "[]"
	inputs[KeyQuestion] = text
	return c.Run(ctx, inputs)
}

func (c *MapRerankDocumentsChain) Apply(ctx context.Context, inputs []map[string]string) ([]string, error) {
	return chains.ChainApplyAsync(c, ctx, inputs)
}

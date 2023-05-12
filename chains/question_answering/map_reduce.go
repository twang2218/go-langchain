package question_answering

import (
	"context"
	"langchain/chains"
	"langchain/llms"
	"langchain/prompts"

	"github.com/sirupsen/logrus"
)

const (
	QuestionPromptTemplate = `Use the following portion of a long document to see if any of the text is relevant to answer the question. 
	Return any relevant text verbatim.
	{` + KeyContext + `}
	Question: {` + KeyQuestion + `}
	Relevant text, if any, or only 'None' if nothing relavant:`
	QuestionChatPromptTemplate = `Use the following portion of a long document to see if any of the text is relevant to answer the question. 
	Return any relevant text verbatim.
	______________________
	{` + KeyContext + `}`
	CombinePromptTemplate = `Given the following extracted parts of a long document and a question, create a final answer. 
	If you don't know the answer, just say that you don't know. Don't try to make up an answer.
	
	QUESTION: Which state/country's law governs the interpretation of the contract?
	=========
	Content: This Agreement is governed by English law and the parties submit to the exclusive jurisdiction of the English courts in  relation to any dispute (contractual or non-contractual) concerning this Agreement save that either party may apply to any court for an  injunction or other relief to protect its Intellectual Property Rights.
	
	Content: No Waiver. Failure or delay in exercising any right or remedy under this Agreement shall not constitute a waiver of such (or any other)  right or remedy.\n\n11.7 Severability. The invalidity, illegality or unenforceability of any term (or part of a term) of this Agreement shall not affect the continuation  in force of the remainder of the term (if any) and this Agreement.\n\n11.8 No Agency. Except as expressly stated otherwise, nothing in this Agreement shall create an agency, partnership or joint venture of any  kind between the parties.\n\n11.9 No Third-Party Beneficiaries.
	
	Content: (b) if Google believes, in good faith, that the Distributor has violated or caused Google to violate any Anti-Bribery Laws (as  defined in Clause 8.5) or that such a violation is reasonably likely to occur,
	=========
	FINAL ANSWER: This Agreement is governed by English law.
	
	QUESTION: What did the president say about Michael Jackson?
	=========
	Content: Madam Speaker, Madam Vice President, our First Lady and Second Gentleman. Members of Congress and the Cabinet. Justices of the Supreme Court. My fellow Americans.  \n\nLast year COVID-19 kept us apart. This year we are finally together again. \n\nTonight, we meet as Democrats Republicans and Independents. But most importantly as Americans. \n\nWith a duty to one another to the American people to the Constitution. \n\nAnd with an unwavering resolve that freedom will always triumph over tyranny. \n\nSix days ago, Russia’s Vladimir Putin sought to shake the foundations of the free world thinking he could make it bend to his menacing ways. But he badly miscalculated. \n\nHe thought he could roll into Ukraine and the world would roll over. Instead he met a wall of strength he never imagined. \n\nHe met the Ukrainian people. \n\nFrom President Zelenskyy to every Ukrainian, their fearlessness, their courage, their determination, inspires the world. \n\nGroups of citizens blocking tanks with their bodies. Everyone from students to retirees teachers turned soldiers defending their homeland.
	
	Content: And we won’t stop. \n\nWe have lost so much to COVID-19. Time with one another. And worst of all, so much loss of life. \n\nLet’s use this moment to reset. Let’s stop looking at COVID-19 as a partisan dividing line and see it for what it is: A God-awful disease.  \n\nLet’s stop seeing each other as enemies, and start seeing each other for who we really are: Fellow Americans.  \n\nWe can’t change how divided we’ve been. But we can change how we move forward—on COVID-19 and other issues we must face together. \n\nI recently visited the New York City Police Department days after the funerals of Officer Wilbert Mora and his partner, Officer Jason Rivera. \n\nThey were responding to a 9-1-1 call when a man shot and killed them with a stolen gun. \n\nOfficer Mora was 27 years old. \n\nOfficer Rivera was 22. \n\nBoth Dominican Americans who’d grown up on the same streets they later chose to patrol as police officers. \n\nI spoke with their families and told them that we are forever in debt for their sacrifice, and we will carry on their mission to restore the trust and safety every community deserves.
		
	Content: More support for patients and families. \n\nTo get there, I call on Congress to fund ARPA-H, the Advanced Research Projects Agency for Health. \n\nIt’s based on DARPA—the Defense Department project that led to the Internet, GPS, and so much more.  \n\nARPA-H will have a singular purpose—to drive breakthroughs in cancer, Alzheimer’s, diabetes, and more. \n\nA unity agenda for the nation. \n\nWe can do this. \n\nMy fellow Americans—tonight , we have gathered in a sacred space—the citadel of our democracy. \n\nIn this Capitol, generation after generation, Americans have debated great questions amid great strife, and have done great things. \n\nWe have fought for freedom, expanded liberty, defeated totalitarianism and terror. \n\nAnd built the strongest, freest, and most prosperous nation the world has ever known. \n\nNow is the hour. \n\nOur moment of responsibility. \n\nOur test of resolve and conscience, of history itself. \n\nIt is in this moment that our character is formed. Our purpose is found. Our future is forged. \n\nWell I know this nation.
	=========
	FINAL ANSWER: The president did not mention Michael Jackson.
	
	QUESTION: {` + KeyQuestion + `}
	=========
	Content: {` + KeyContext + `}
	=========
	FINAL ANSWER:`
	CombineChatPromptTemplate = `Given the following extracted parts of a long document and a question, create a final answer. 
	If you don't know the answer, just say that you don't know. Don't try to make up an answer.
	Content: {` + KeyContext + `}
	=========
	Question: {` + KeyQuestion + `}
	=========
	FINAL ANSWER:`
)

type MapReduceDocumentsChain struct {
	MapChain    chains.Chain
	ReduceChain chains.Chain
}

func NewMapReduceDocumentsChain(llm llms.LLM, questionPrompt, combinePrompt prompts.Template) *MapReduceDocumentsChain {
	if questionPrompt == nil {
		questionPrompt = prompts.NewPromptTemplateByTemplate(QuestionPromptTemplate)
	}
	if combinePrompt == nil {
		combinePrompt = prompts.NewPromptTemplateByTemplate(CombineChatPromptTemplate)
	}
	return &MapReduceDocumentsChain{
		MapChain:    chains.NewLLMChain(llm, questionPrompt),
		ReduceChain: NewStuffDocumentsChain(llm, combinePrompt),
	}
}

func (c *MapReduceDocumentsChain) Run(ctx context.Context, inputs map[string]string) (string, error) {
	logrus.Tracef("MapReduceDocumentsChain.Run(): %v", inputs)
	// Parse documents
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

func (c *MapReduceDocumentsChain) CombineDocs(ctx context.Context, docs []string, inputs map[string]string) (string, error) {
	//  Map each document to a question
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

	results, err := c.MapChain.Apply(ctx, docs_inputs)
	if err != nil {
		return "", err
	}

	//  Reduce the results to a single answer
	err = PutDocsToInputs(results, inputs, KeyDocuments)
	if err != nil {
		return "", err
	}

	return c.ReduceChain.Run(ctx, inputs)
}

func (c *MapReduceDocumentsChain) RunText(ctx context.Context, input string) (string, error) {
	inputs := make(map[string]string)
	inputs[KeyDocuments] = "[]"
	inputs[KeyQuestion] = input
	return c.Run(ctx, inputs)
}

func (c *MapReduceDocumentsChain) Apply(ctx context.Context, inputs []map[string]string) ([]string, error) {
	return chains.ChainApplyAsync(c, ctx, inputs)
}

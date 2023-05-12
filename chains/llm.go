package chains

import (
	"context"
	"fmt"
	"langchain/llms"
	"langchain/prompts"

	"github.com/sirupsen/logrus"
)

type LLMChain struct {
	LLM    llms.LLM
	Prompt prompts.Template
}

func NewLLMChain(llm llms.LLM, prompt prompts.Template) *LLMChain {
	if prompt == nil {
		prompt = prompts.NewPromptTemplateByTemplate("{prompt}")
	}
	return &LLMChain{
		LLM:    llm,
		Prompt: prompt,
	}
}

func (c *LLMChain) Run(ctx context.Context, inputs map[string]string) (string, error) {
	logrus.Tracef("LLMChain.Run(): %v", inputs)
	prompt := c.Prompt.Format(inputs)
	resp, err := c.LLM.Chat(ctx, prompt)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *LLMChain) RunText(ctx context.Context, input string) (string, error) {
	inputs := map[string]string{}
	variables := c.Prompt.GetVariables()
	if len(variables) != 1 {
		return "", fmt.Errorf("LLMChain.Run(): PromptTemplate can only have 1 variable, but now it contains %d variables(%v)", len(variables), variables)
	}
	for _, variable := range variables {
		inputs[variable] = input
	}
	return c.Run(ctx, inputs)
}

func (c *LLMChain) Apply(ctx context.Context, inputs []map[string]string) ([]string, error) {
	return ChainApplyAsync(c, ctx, inputs)
}

package llms

import (
	"context"
	"fmt"
	"langchain/callbacks"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultOpenAIMaxRetries  = 3
	DefaultOpenAIModel       = "gpt-3.5-turbo"
	DefaultOpenAITemperature = 0.1
	DefaultOpenAITopP        = 1.0
	DefaultOpenAIN           = 1
)

type OpenAI struct {
	OpenAIKey        string
	Model            string
	MaxTokens        int
	Temperature      float32
	TopP             float32
	N                int
	Stop             []string
	PresencePenalty  float32
	FrequencyPenalty float32
	Stream           bool
	MaxRetries       int
	Callbacks        *callbacks.CallbackManager
	client           *openai.Client
}

func NewOpenAI(openAIKey string) *OpenAI {
	l := &OpenAI{OpenAIKey: openAIKey}
	l.init()
	return l
}

func (l *OpenAI) init() error {
	if l.OpenAIKey == "" {
		l.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	}
	if l.OpenAIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set")
	}
	if l.client == nil {
		l.client = openai.NewClient(l.OpenAIKey)
	}
	if l.Model == "" {
		l.Model = DefaultOpenAIModel
	}
	if l.Temperature < 0 {
		l.Temperature = DefaultOpenAITemperature
	}
	if l.TopP < 0 {
		l.TopP = DefaultOpenAITopP
	}
	if l.N < 0 {
		l.N = DefaultOpenAIN
	}
	if l.MaxRetries <= 0 {
		l.MaxRetries = DefaultOpenAIMaxRetries
	}
	return nil
}

func (l *OpenAI) Chat(ctx context.Context, prompt string) (string, error) {
	if err := l.init(); err != nil {
		return "", err
	}

	//	对于要求流输出的，使用ChatStream()来处理请求
	if l.Stream {
		return l.ChatStream(ctx, prompt)
	}

	l.Callbacks.OnLLMStart(prompt)

	var err error
	var resp openai.ChatCompletionResponse
	for i := 0; i < l.MaxRetries; i++ {
		resp, err = l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			Temperature:      l.Temperature,
			TopP:             l.TopP,
			N:                l.N,
			Stop:             l.Stop,
			PresencePenalty:  l.PresencePenalty,
			FrequencyPenalty: l.FrequencyPenalty,
		})
		if err == nil {
			break
		}
		// need to retry
		log.Warnf("openai api error: %s", err.Error())
		l.Callbacks.OnLLMError(err)
	}

	if len(resp.Choices) == 0 {
		err = fmt.Errorf("openai api error: response is empty")
		l.Callbacks.OnLLMError(err)
		return "", err
	}
	if len(resp.Choices[0].Message.Content) == 0 {
		err = fmt.Errorf("openai api error: response content is empty")
		l.Callbacks.OnLLMError(err)
		return "", err
	}

	//	TODO: remove this as it can be implemented in the callback
	log.Tracef("prompt: %s", prompt)
	if len(prompt) > 5000 {
		log.Warnf("prompt length (%d) is beyond 5000, it might be truncated and leading to wrong answer.", len(prompt))
	}

	l.Callbacks.OnLLMEnd(resp.Choices[0].Message.Content)
	//	TODO: remove this as it can be implemented in the callback
	log.Debugf("openai api response: %+v", resp)
	return resp.Choices[0].Message.Content, nil
}

func (l *OpenAI) ChatStream(ctx context.Context, prompt string) (string, error) {
	if err := l.init(); err != nil {
		return "", err
	}

	l.Callbacks.OnLLMStart(prompt)

	var err error
	var stream *openai.ChatCompletionStream
	for i := 0; i < l.MaxRetries; i++ {
		stream, err = l.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			Temperature:      l.Temperature,
			TopP:             l.TopP,
			N:                l.N,
			Stop:             l.Stop,
			PresencePenalty:  l.PresencePenalty,
			FrequencyPenalty: l.FrequencyPenalty,
			Stream:           true,
		})
		if err == nil {
			break
		}
		// need to retry
		log.Warnf("openai api error: %s", err.Error())
		l.Callbacks.OnLLMError(err)
	}

	defer stream.Close()

	var sb strings.Builder
	for {
		resp, err := stream.Recv()
		if err != nil {
			l.Callbacks.OnLLMError(err)
			return sb.String(), err
		}
		if len(resp.Choices) == 0 || resp.Choices[0].FinishReason == "stop" {
			l.Callbacks.OnLLMNewToken("[DONE]")
			return sb.String(), nil
		}
		content := resp.Choices[0].Delta.Content
		sb.WriteString(content)
		l.Callbacks.OnLLMNewToken(content)
	}
}

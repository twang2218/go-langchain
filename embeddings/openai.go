package embeddings

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	OpenAIKey string
	Model     string
	client    *openai.Client
}

func NewOpenAI(openAIKey string) *OpenAI {
	return &OpenAI{
		OpenAIKey: openAIKey,
		Model:     "text-embedding-ada-002",
		client:    openai.NewClient(openAIKey),
	}
}

func (l *OpenAI) Embedding(ctx context.Context, text string) ([]float32, error) {
	var em openai.EmbeddingModel
	err := em.UnmarshalText([]byte(l.Model))
	if err != nil {
		return nil, err
	}

	resp, err := l.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: em,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("embedding is empty")
	}

	return resp.Data[0].Embedding, nil
}

func (l *OpenAI) Embeddings(ctx context.Context, texts []string) ([][]float32, error) {
	var em openai.EmbeddingModel
	err := em.UnmarshalText([]byte(l.Model))
	if err != nil {
		return nil, err
	}

	resp, err := l.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: em,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("embedding is empty")
	}

	var embeddings [][]float32
	for _, d := range resp.Data {
		embeddings = append(embeddings, d.Embedding)
	}

	return embeddings, nil
}

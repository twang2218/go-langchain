package embeddings

import "context"

type Embedding interface {
	// Embedding 获取文本的向量表示
	Embedding(ctx context.Context, text string) ([]float32, error)
	Embeddings(ctx context.Context, texts []string) ([][]float32, error)
}

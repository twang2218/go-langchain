package vectorstores

import (
	"context"
	"langchain/docstore"
)

type VectorStore interface {
	AddDocuments(ctx context.Context, docs []docstore.Document) error
	SimilaritySearch(ctx context.Context, query string, limit int) ([]docstore.Document, error)
	SimilaritySearchWithScore(ctx context.Context, query string, limit int) ([]docstore.Document, error)
	Clean(ctx context.Context) error
	Dump(ctx context.Context) ([]docstore.Document, error)
}

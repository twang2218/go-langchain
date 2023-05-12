package vectorstores

import (
	"context"
	"fmt"
	"langchain/docstore"
	"langchain/embeddings"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type Weaviate struct {
	URL       string
	IndexName string
	Embedding embeddings.Embedding
	Headers   map[string]string
	client    *weaviate.Client
}

func (w *Weaviate) init(ctx context.Context) error {
	//	check client
	if w.client == nil {
		u, err := url.Parse(w.URL)
		if err != nil {
			return err
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		cfg := weaviate.Config{
			Scheme:  u.Scheme,
			Host:    u.Host,
			Headers: w.Headers,
		}
		w.client, err = weaviate.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("Weaviate.init(): create client error: %v", getError(err))
		}
	}

	//	check index
	schema, err := w.client.Schema().ClassGetter().WithClassName(w.IndexName).Do(ctx)
	if err != nil || schema == nil {
		//	create index
		schema = &models.Class{
			Class: w.IndexName,
			Properties: []*models.Property{
				{
					Name:     "text",
					DataType: []string{"text"},
				},
				// metadata schema will be generated automatically during the data insertion
			},
		}
		if err := w.client.Schema().ClassCreator().WithClass(schema).Do(ctx); err != nil {
			return fmt.Errorf("Weaviate.init(): create schema error: %v", getError(err))
		}

	}

	return nil
}

func (w *Weaviate) AddDocuments(ctx context.Context, docs []docstore.Document) error {
	if err := w.init(ctx); err != nil {
		return err
	}

	objs := make([]*models.Object, 0, len(docs))
	hasEmbedding := false
	missingEmbedding := false
	for _, doc := range docs {
		if doc.Embedding == nil {
			missingEmbedding = true
		} else {
			hasEmbedding = true
		}

		//	add document
		props := map[string]interface{}{
			"text": doc.Content,
		}
		for key, value := range doc.Metadata {
			props[key] = value
		}
		objs = append(objs, &models.Object{
			Class:      w.IndexName,
			Properties: props,
			Vector:     doc.Embedding,
		})
	}
	//	检查计算 Embedding
	if missingEmbedding && hasEmbedding {
		//	说明有些文档没有预先计算好的embedding，有些文档有预先计算好的embedding
		//	这种情况下，不批量处理，而是逐个处理
		for i, obj := range objs {
			if obj.Vector == nil {
				embedding, err := w.Embedding.Embedding(ctx, docs[i].Content)
				if err != nil {
					return err
				}
				obj.Vector = embedding
			}
		}
	} else if missingEmbedding {
		//	说明所有文档都没有预先计算好的embedding
		//	这种情况下，批量处理。为了防止文档过多，导致批量处理失败，这里限制每次批量处理的文档数量
		const BATCH_EMBEDDING_SIZE = 50
		texts := make([]string, 0, len(objs))
		for _, doc := range docs {
			texts = append(texts, doc.Content)
		}
		for i := 0; i < len(objs); i += BATCH_EMBEDDING_SIZE {
			end := i + BATCH_EMBEDDING_SIZE
			if end > len(objs) {
				end = len(objs)
			}
			embeddings, err := w.Embedding.Embeddings(ctx, texts[i:end])
			if err != nil {
				return err
			}
			for j, embedding := range embeddings {
				objs[i+j].Vector = embedding
			}
		}
	}

	_, err := w.client.Batch().ObjectsBatcher().
		WithObjects(objs...).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("Weaviate.AddDocuments(): batch add documents error: %v", getError(err))
	}

	return nil
}

func (w *Weaviate) SimilaritySearch(ctx context.Context, query string, limit int) ([]docstore.Document, error) {
	if err := w.init(ctx); err != nil {
		return nil, err
	}

	//	embedding
	embedding, err := w.Embedding.Embedding(ctx, query)
	if err != nil {
		return nil, err
	}

	//	get metadata keys
	schema, err := w.client.Schema().ClassGetter().WithClassName(w.IndexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("Weaviate.SimilaritySearch(): get schema error: %v", getError(err))
	}
	metadataKeys := make([]string, 0)
	for _, prop := range schema.Properties {
		if prop.Name != "text" {
			metadataKeys = append(metadataKeys, prop.Name)
		}
	}

	//	fields
	fields := []graphql.Field{
		{Name: "text"},
	}
	for _, key := range metadataKeys {
		fields = append(fields, graphql.Field{Name: key})
	}
	//	search
	resp, err := w.client.GraphQL().Get().
		WithClassName(w.IndexName).
		WithFields(fields...).
		WithNearVector(
			w.client.GraphQL().NearVectorArgBuilder().WithVector(embedding),
		).
		WithLimit(limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("Weaviate.SimilaritySearch(): search error: %v", getError(err))
	}

	//	parse result
	docs := make([]docstore.Document, 0)
	if gets, ok := resp.Data["Get"].(map[string]interface{}); ok {
		if objects, ok := gets[w.IndexName].([]interface{}); ok {
			for _, object := range objects {
				if obj, ok := object.(map[string]interface{}); ok {
					metadata := make(map[string]string)
					for _, key := range metadataKeys {
						metadata[key] = obj[key].(string)
					}
					docs = append(docs, docstore.Document{
						Content:  obj["text"].(string),
						Metadata: metadata,
					})
				}
			}
		}
	}
	return docs, nil
}

func (w *Weaviate) SimilaritySearchWithScore(ctx context.Context, query string, limit int) ([]docstore.Document, error) {
	if err := w.init(ctx); err != nil {
		return nil, err
	}

	//	get metadata keys
	schema, err := w.client.Schema().ClassGetter().WithClassName(w.IndexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("Weaviate.SimilaritySearch(): get schema error: %v", getError(err))
	}
	metadataKeys := make([]string, 0)
	for _, prop := range schema.Properties {
		if prop.Name != "text" {
			metadataKeys = append(metadataKeys, prop.Name)
		}
	}

	//	fields
	fields := []graphql.Field{
		{Name: "text"},
	}
	for _, key := range metadataKeys {
		fields = append(fields, graphql.Field{Name: key})
	}
	//	fields for the score
	fields = append(fields, graphql.Field{
		Name: "_additional",
		Fields: []graphql.Field{
			{Name: "certainty"},
		},
	})
	//	search
	//	embedding
	embedding, err := w.Embedding.Embedding(ctx, query)
	if err != nil {
		return nil, err
	}
	vec := w.client.GraphQL().NearVectorArgBuilder().WithVector(embedding)
	resp, err := w.client.GraphQL().Get().
		WithClassName(w.IndexName).
		WithFields(fields...).
		WithNearVector(vec).
		WithLimit(limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("Weaviate.SimilaritySearch(): search error: %v", getError(err))
	}

	//	parse result
	docs := make([]docstore.Document, 0)
	if gets, ok := resp.Data["Get"].(map[string]interface{}); ok {
		if objects, ok := gets[w.IndexName].([]interface{}); ok {
			for _, object := range objects {
				if obj, ok := object.(map[string]interface{}); ok {
					metadata := make(map[string]string)
					for _, key := range metadataKeys {
						metadata[key] = obj[key].(string)
					}
					metadata["score"] = fmt.Sprintf("%v", obj["_additional"].(map[string]interface{})["certainty"].(float64))
					docs = append(docs, docstore.Document{
						Content:  obj["text"].(string),
						Metadata: metadata,
					})
				}
			}
		}
	}
	return docs, nil
}

func (w *Weaviate) Clean(ctx context.Context) error {
	if err := w.init(ctx); err != nil {
		return err
	}

	//	delete all objects
	err := w.client.Schema().ClassDeleter().WithClassName(w.IndexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("Weaviate.Clean(): delete all objects error: %v", getError(err))
	}
	return nil
}

func (w *Weaviate) Dump(ctx context.Context) ([]docstore.Document, error) {
	if err := w.init(ctx); err != nil {
		return nil, err
	}

	//	get metadata keys
	schema, err := w.client.Schema().ClassGetter().WithClassName(w.IndexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("Weaviate.Dump(): get schema error: %v", getError(err))
	}
	metadataKeys := make([]string, 0)
	for _, prop := range schema.Properties {
		if prop.Name != "text" {
			metadataKeys = append(metadataKeys, prop.Name)
		}
	}

	//	search
	resp, err := w.client.Data().ObjectsGetter().
		WithClassName(w.IndexName).
		WithAdditional("vector").
		WithLimit(10000).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("Weaviate.Dump(): search error: %v", getError(err))
	}

	//	parse result
	docs := make([]docstore.Document, 0)
	for _, object := range resp {
		props := object.Properties.(map[string]interface{})
		metadata := make(map[string]string)
		for _, key := range metadataKeys {
			if value, ok := props[key]; ok {
				metadata[key] = value.(string)
			}
		}
		docs = append(docs, docstore.Document{
			Content:   props["text"].(string),
			Metadata:  metadata,
			Embedding: object.Vector,
		})
	}
	return docs, nil
}

// getError 获取错误
func getError(err error) error {
	if err, ok := err.(*fault.WeaviateClientError); ok {
		if err.StatusCode == -1 {
			switch err.DerivedFromError.(type) {
			case *fault.WeaviateClientError:
				return getError(err.DerivedFromError)
			case *url.Error:
				return fmt.Errorf("%v", err.DerivedFromError.(*url.Error).Error())
			default:
				return fmt.Errorf("%v", err.DerivedFromError.Error())
			}
		}
		return fmt.Errorf("[%d] %v", err.StatusCode, err.Msg)
	}
	return err
}

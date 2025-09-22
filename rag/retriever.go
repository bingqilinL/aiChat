package rag

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
)

func newRetriever(ctx context.Context, embedder *ark.Embedder, client cli.Client) *milvus.Retriever {
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      client,
		Embedding:   embedder,
		Collection:  "test",
		VectorField: "vector",
		TopK:        1,
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
	})

	if err != nil {
		panic(err)
	}
	return retriever
}

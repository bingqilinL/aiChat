package rag

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"os"
)

func newEmbedder(ctx context.Context) *ark.Embedder {
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		Model:  os.Getenv("EMBEDDER"),
		APIKey: os.Getenv("ARK_API_KEY"),
	})
	if err != nil {
		panic(err)
	}
	return embedder
}

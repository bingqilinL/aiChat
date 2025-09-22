package rag

import (
	"context"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/document"
)

func newMarkdown(ctx context.Context) document.Transformer {
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}
	return splitter
}
func newSemantic(ctx context.Context, embedder *ark.Embedder) document.Transformer {
	splitter, err := semantic.NewSplitter(ctx, &semantic.Config{
		MinChunkSize: 1024,
		Embedding:    embedder,
		BufferSize:   5,
		Percentile:   0.95,
	})
	if err != nil {
		panic(err)
	}
	return splitter
}

func newRecursive(ctx context.Context) document.Transformer {
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   1000,
		OverlapSize: 200,
		//Separators:  []string{"\n\n", "\n", "。", "！", "？"},
		KeepType: recursive.KeepTypeEnd,
	})
	if err != nil {
		panic(err)
	}
	return splitter
}

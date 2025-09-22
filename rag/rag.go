package rag

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	retriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
	"os"
)

type Rag struct {
	embedder  *ark.Embedder
	indexer   *milvus.Indexer
	retriever *retriever.Retriever
	splitter  document.Transformer
}

func NewRAG(ctx context.Context) (*Rag, error) {
	embedder := newEmbedder(ctx)
	client := initClient(ctx)
	indexer := newIndexer(ctx, embedder, client)
	retriever := newRetriever(ctx, embedder, client)
	//splitter := newMarkdown(ctx)
	//splitter := newSemantic(ctx, embedder)
	splitter := newRecursive(ctx)
	return &Rag{
		embedder:  embedder,
		indexer:   indexer,
		retriever: retriever,
		splitter:  splitter,
	}, nil
}
func initClient(ctx context.Context) cli.Client {
	client, err := cli.NewClient(ctx, cli.Config{
		Address: "localhost:19530",
		DBName:  "MyEino",
	})
	if err != nil {
		panic(err)
	}
	return client
}
func loadDoc(path string) ([]*schema.Document, error) {
	ct, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	docs := make([]*schema.Document, len(ct))
	for i, v := range ct {
		doc, err := os.ReadFile(path + "/" + v.Name())
		if err != nil {
			return nil, err
		}
		docs[i] = &schema.Document{
			ID:      v.Name(),
			Content: string(doc),
		}
	}
	return docs, nil
}

func (r *Rag) Store(ctx context.Context, path string) error {
	docs, err := loadDoc(path)
	if err != nil {
		panic(err)
	}

	result, err := r.splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	_, err = r.indexer.Store(ctx, result)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Rag) StoreDoc(ctx context.Context, id string, content string) error {
	docs := []*schema.Document{
		{
			ID:      id,
			Content: content,
		},
	}

	result, err := r.splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	_, err = r.indexer.Store(ctx, result)
	if err != nil {
		panic(err)
	}
	return nil
}
func (r *Rag) Find(ctx context.Context, question string) []*schema.Document {
	res, err := r.retriever.Retrieve(ctx, question)
	if err != nil {
		panic(err)
	}
	return res
}

package main

import (
	"context"
	"fmt"
	"github.com/ryomak/llm-qa-go-example/langchain"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	weaviateClient := weaviate.New(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	})
	ctx := context.Background()
	if err := weaviateClient.Schema().ClassCreator().WithClass(&models.Class{
		Class:       langchain.WeaviateIndexName,
		Description: "qa class",
		VectorIndexConfig: map[string]any{
			"distance": "cosine",
		},
		ModuleConfig: map[string]any{},
		Properties: []*models.Property{
			{
				Name:        langchain.WeaviatePropertyTextName,
				Description: "question",
				DataType:    []string{"text"},
			},
			{
				Name:        langchain.WeaviatePropertyNameSpaceName,
				Description: "namespace",
				DataType:    []string{"text"},
			},
			{
				Name:        langchain.WeaviatePropertyKind,
				Description: "kind",
				DataType:    []string{"text"},
			},
		},
	}).Do(ctx); err != nil {
		panic(err)
	}
	fmt.Println("created")
}

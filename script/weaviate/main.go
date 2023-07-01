package main

import (
	"context"
	"fmt"
	"github.com/ryomak/langchain-go-example/qa"
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
		Class:       qa.WeaviateIndexName,
		Description: "qa class",
		VectorIndexConfig: map[string]any{
			"distance": "cosine",
		},
		ModuleConfig: map[string]any{},
		Properties: []*models.Property{
			{
				Name:        qa.WeaviatePropertyTextName,
				Description: "document text",
				DataType:    []string{"text"},
			},
			{
				Name:        qa.WeaviatePropertyNameSpaceName,
				Description: "namespace",
				DataType:    []string{"text"},
			},
		},
	}).Do(ctx); err != nil {
		panic(err)
	}
	fmt.Println("created")
}

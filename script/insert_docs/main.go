package main

import (
	"context"
	"fmt"
	"github.com/ryomak/llm-qa-go-example/langchain"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"os"
)

func main() {
	chain, err := langchain.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	filePaths := []string{
		"./langchain/llm.go",
		"./docker-compose.yml",
		"./LICENSE",
		"go.mod",
	}
	for _, filePath := range filePaths {

		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		docs, err := documentloaders.NewText(file).LoadAndSplit(
			context.Background(),
			textsplitter.NewRecursiveCharacter(),
		)
		if err != nil {
			panic(err)
		}

		for _, v := range docs {
			if err := chain.AddDocument(ctx, "default", v.PageContent); err != nil {
				panic(err)
			}
		}

	}
	fmt.Println("done")
}

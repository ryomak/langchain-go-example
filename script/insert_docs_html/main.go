package main

import (
	"context"
	"fmt"
	"github.com/ryomak/langchain-go-example/qa"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"os"
)

func main() {
	chain, err := qa.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	file, err := os.Open("./script/insert_docs_html/qa.html")
	if err != nil {
		panic(err)
	}
	loader := documentloaders.NewHTML(file)
	docs, err := loader.LoadAndSplit(
		context.Background(),
		textsplitter.RecursiveCharacter{
			Separators:   []string{"\n\n", "\n", " ", ""},
			ChunkSize:    800,
			ChunkOverlap: 200,
		},
	)
	if err != nil {
		panic(err)
	}
	for _, v := range docs {
		if err := chain.AddDocument(ctx, qa.NameSpaceHTML, v.PageContent); err != nil {
			panic(err)
		}
	}
	fmt.Println("done")
}

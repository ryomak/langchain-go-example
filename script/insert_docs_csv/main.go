package main

import (
	"context"
	"fmt"
	"github.com/ryomak/langchain-go-example/qa"
	"github.com/tmc/langchaingo/documentloaders"
	"os"
)

func main() {
	chain, err := qa.New()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	file, err := os.Open("./script/insert_docs_csv/qa.csv")
	if err != nil {
		panic(err)
	}
	loader := documentloaders.NewCSV(file)
	docs, err := loader.Load(context.Background())
	if err != nil {
		panic(err)
	}
	for _, v := range docs {
		if err := chain.AddDocument(ctx, qa.NameSpaceCSV, v.PageContent); err != nil {
			panic(err)
		}
	}
	fmt.Println("done")
}

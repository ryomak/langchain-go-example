package main

import (
	"context"
	"github.com/ryomak/llm-qa-go-example/langchain"
)

func main() {
	llm, err := langchain.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	question := "メロスの家族は？"
	result, err := llm.Run(ctx, question)
	if err != nil {
		panic(err)
	}
	println(result)
}

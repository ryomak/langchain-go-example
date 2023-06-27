package main

import (
	"context"
	"fmt"
	"github.com/ryomak/llm-qa-go-example/langchain"
)

func main() {
	llm, err := langchain.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	question := "Nissyとはどういう人ですか?"
	result, err := llm.QA(ctx, question)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

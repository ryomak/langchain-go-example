package main

import (
	"context"
	"fmt"
	"github.com/ryomak/llm-qa-go-example/qa"
)

func main() {
	qaBot, err := qa.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	question := "Nissyの誕生日は?"
	result, err := qaBot.Answer(ctx, question)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

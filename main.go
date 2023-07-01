package main

import (
	"context"
	"fmt"
	"github.com/ryomak/langchain-go-example/qa"
)

func main() {
	qaBot, err := qa.New()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	for _, v := range []struct {
		question  string
		nameSpace qa.NameSpace
	}{
		{
			question:  "Nissyの生年月日はなんですか?また、Nissyの本名はなんですか？",
			nameSpace: qa.NameSpaceCSV,
		},
		{
			question:  "Go1.21に追加された3つのbuilt-insはなんですか？",
			nameSpace: qa.NameSpaceHTML,
		},
	} {
		result, err := qaBot.Answer(ctx, v.nameSpace, v.question)
		if err != nil {
			panic(err)
		}

		fmt.Println("=====================================")
		fmt.Println("kind:\n", v.nameSpace)
		fmt.Println("question:\n", v.question)
		fmt.Println("result:\n", result)
		fmt.Println("=====================================")
	}
}

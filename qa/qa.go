package qa

import (
	"context"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

const (
	WeaviateIndexName = "QA_2023"

	WeaviatePropertyTextName      = "text"
	WeaviatePropertyNameSpaceName = "namespace"
)

type NameSpace string

const (
	NameSpaceCSV  NameSpace = "csv"
	NameSpaceHTML NameSpace = "html"
)

type QA struct {
	llm   llms.LanguageModel
	store vectorstores.VectorStore
}

func New() (*QA, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, err
	}
	e, err := embeddings.NewOpenAI()
	if err != nil {
		return nil, err
	}
	store, err := weaviate.New(
		weaviate.WithScheme("http"),
		weaviate.WithHost("localhost:8080"),
		weaviate.WithEmbedder(e),
		weaviate.WithIndexName(WeaviateIndexName),
		weaviate.WithTextKey(WeaviatePropertyTextName),
		weaviate.WithNameSpaceKey(WeaviatePropertyNameSpaceName),
	)
	if err != nil {
		return nil, err
	}
	return &QA{
		llm:   llm,
		store: store,
	}, nil
}

func (l *QA) AddDocument(ctx context.Context, namespace NameSpace, content string) error {
	return l.store.AddDocuments(ctx, []schema.Document{
		{
			PageContent: content,
		},
	}, vectorstores.WithNameSpace(string(namespace)))
}

func (l *QA) Answer(ctx context.Context, namespace NameSpace, question string) (string, error) {
	prompt := prompts.NewPromptTemplate(
		`## Introduction 
あなたはカスタマーサポートです。丁寧な回答を心がけてください。
以下のContextを使用して、日本語で質問に答えてください。Contextから答えがわからない場合は、「わかりません」と回答してください。

## 質問
{{.question}}

## Context
{{.context}}

日本語での回答:`,

		[]string{"context", "question"},
	)

	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(l.llm, prompt))

	result, err := chains.Run(
		ctx,
		chains.NewRetrievalQA(
			combineChain,
			vectorstores.ToRetriever(
				l.store,
				5,
				vectorstores.WithNameSpace(string(namespace)),
			),
		),
		question,
		chains.WithModel("gpt-4-0613"),
	)

	if err != nil {
		return "", err
	}
	return result, nil
}

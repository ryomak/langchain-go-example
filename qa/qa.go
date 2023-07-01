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
	WeaviateIndexName = "QA_20230701_2"

	WeaviatePropertyTextName      = "text"
	WeaviatePropertyNameSpaceName = "namespace"
)

type qa struct {
	llm   llms.LanguageModel
	store vectorstores.VectorStore
}

func New() (*qa, error) {
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
	return &qa{
		llm:   llm,
		store: store,
	}, nil
}

func (l *qa) AddDocument(ctx context.Context, content string) error {
	return l.store.AddDocuments(ctx, []schema.Document{
		{
			PageContent: content,
		},
	})
}

func (l *qa) Answer(ctx context.Context, question string) (string, error) {
	prompt := prompts.NewPromptTemplate(
		`## 依頼
あなたはカスタマーサポートです。丁寧な回答を心がけてください。
以下の過去の質問と回答を使用して,質問に答えてください。
過去の質問と回答から答えがわからない場合は、「わかりません」と回答してください。

## 質問
{{.question}}

## 過去の質問と回答
{{.context}}
`,

		[]string{"context", "question"},
	)

	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(l.llm, prompt))

	result, err := chains.Run(
		ctx,
		chains.NewRetrievalQA(
			combineChain,
			vectorstores.ToRetriever(l.store, 3),
		),
		question,
		chains.WithModel("gpt-4"),
	)
	if err != nil {
		return "", err
	}
	return result, nil
}

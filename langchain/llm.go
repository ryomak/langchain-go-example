package langchain

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
	WeaviateIndexName = "QA_20230626_1"

	WeaviatePropertyTextName      = "text"
	WeaviatePropertyNameSpaceName = "namespace"
	WeaviatePropertyKind          = "kind"
)

type chain struct {
	llm   llms.LLM
	store vectorstores.VectorStore
}

func New() (*chain, error) {
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
		weaviate.WithQueryAttrs([]string{WeaviatePropertyKind}),
		weaviate.WithTextKey(WeaviatePropertyTextName),
		weaviate.WithNameSpaceKey(WeaviatePropertyNameSpaceName),
	)
	if err != nil {
		return nil, err
	}
	return &chain{
		llm:   llm,
		store: store,
	}, nil
}

func (l *chain) AddDocument(ctx context.Context, kind string, content string) error {
	return l.store.AddDocuments(ctx, []schema.Document{
		{
			PageContent: content,
			Metadata: map[string]any{
				WeaviatePropertyKind: kind,
			},
		},
	})
}

func (l *chain) QA(ctx context.Context, question string) (string, error) {
	prompt := prompts.NewPromptTemplate(
		`あなたはカスタマーサポートです。丁寧な回答を心がけてください。
以下のContextを使用して質問に答えてください。
Contextから答えがわからない場合は、「わかりません」と回答してください。

Context:
{{.context}}
質問:
{{.question}}
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

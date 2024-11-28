package rag

import (
	"context"
	"fmt"
	"github.com/konradkrasno/ragserver/config"
	"github.com/konradkrasno/ragserver/models"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
	"strings"
)

type Rag struct {
	Ctx       context.Context
	Config    *config.Config
	WvStore   weaviate.Store
	LLMClient llms.Model
}

func New() (*Rag, error) {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	ollamaClient, err := ollama.New(
		ollama.WithServerURL(fmt.Sprintf("%s:%s", cfg.OllamaUrl, cfg.OllamaPort)),
		ollama.WithModel(cfg.LLM),
	)
	if err != nil {
		return nil, err
	}

	emb, err := embeddings.NewEmbedder(ollamaClient)
	if err != nil {
		return nil, err
	}

	wvStore, err := weaviate.New(
		weaviate.WithEmbedder(emb),
		weaviate.WithScheme(cfg.Scheme),
		weaviate.WithHost(fmt.Sprintf("%s:%s", cfg.WvHost, cfg.WvPort)),
		weaviate.WithIndexName(cfg.IndexName),
	)
	if err != nil {
		return nil, err
	}

	return &Rag{
		Ctx:       ctx,
		Config:    cfg,
		WvStore:   wvStore,
		LLMClient: ollamaClient,
	}, nil
}

func (rs *Rag) AddDocuments(adr models.AddDocumentsRequest) error {
	wvDocs := make([]schema.Document, len(adr.Documents))
	for i, doc := range adr.Documents {
		wvDocs[i] = schema.Document{
			PageContent: doc.Text,
		}
	}

	_, err := rs.WvStore.AddDocuments(rs.Ctx, wvDocs)
	if err != nil {
		return err
	}

	return nil
}

func (rs *Rag) query(qr models.QueryRequest) (string, error) {
	docs, err := rs.WvStore.SimilaritySearch(rs.Ctx, qr.Content, rs.Config.DocumentsRetrievalNumber)
	if err != nil {
		return "", err
	}
	docContents := make([]string, len(docs))
	for i, doc := range docs {
		docContents[i] = doc.PageContent
	}

	ragQuery := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(docContents, "\n"))
	return llms.GenerateFromSinglePrompt(rs.Ctx, rs.LLMClient, ragQuery, llms.WithModel(rs.Config.LLM))
}

func (rs *Rag) Query(qr models.QueryRequest) {
	answer, err := rs.query(qr)

	if err != nil {
		fmt.Println("error occurred:", err)
	}
	fmt.Println("answer:", answer)

	// todo: save results in db
	// todo: send results to Websocket
}

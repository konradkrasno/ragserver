package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/konradkrasno/ragserver/broker"
	"github.com/konradkrasno/ragserver/environment"
	"github.com/konradkrasno/ragserver/models"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
	"log"
	"strings"
)

type Rag struct {
	Ctx       context.Context
	Envs      *environment.Envs
	WvStore   weaviate.Store
	LLMClient llms.Model
	Broker    broker.Broker
}

func New(envs *environment.Envs, broker broker.Broker) (*Rag, error) {
	ctx := context.Background()

	ollamaClient, err := ollama.New(
		ollama.WithServerURL(
			fmt.Sprintf(
				"%s:%s", fmt.Sprintf("http://%s", envs.OllamaHost), envs.OllamaPort,
			),
		),
		ollama.WithModel(envs.LLM),
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
		weaviate.WithScheme(envs.WvScheme),
		weaviate.WithHost(fmt.Sprintf("%s:%s", envs.WvHost, envs.WvPort)),
		weaviate.WithIndexName(envs.WvIndexName),
	)
	if err != nil {
		return nil, err
	}

	return &Rag{
		Ctx:       ctx,
		Envs:      envs,
		WvStore:   wvStore,
		LLMClient: ollamaClient,
		Broker:    broker,
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
	docs, err := rs.WvStore.SimilaritySearch(rs.Ctx, qr.Content, rs.Envs.WvDocumentsRetrievalNumber)
	if err != nil {
		return "", err
	}
	docContents := make([]string, len(docs))
	for i, doc := range docs {
		docContents[i] = doc.PageContent
	}

	ragQuery := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(docContents, "\n"))
	return llms.GenerateFromSinglePrompt(rs.Ctx, rs.LLMClient, ragQuery, llms.WithModel(rs.Envs.LLM))
}

func (rs *Rag) Query(qr models.QueryRequest) {
	answer, err := rs.query(qr)
	if err != nil {
		log.Println(err)
		answer = fmt.Sprintf("error occurred: %s", err)
	}

	resp := models.QueryResponse{
		Query:  qr.Content,
		Answer: answer,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	rs.Broker.Publish(rs.Envs.RabbitMQAnswerExchange, qr.SessionId, data)
}

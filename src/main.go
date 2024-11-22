package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/config"
	"github.com/konradkrasno/ragserver/server"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	ollamaClient, err := ollama.New(
		ollama.WithServerURL(fmt.Sprintf("%s:%s", cfg.OllamaHost, cfg.OllamaPort)),
		ollama.WithModel(cfg.LLM),
	)
	if err != nil {
		panic(err)
	}

	emb, err := embeddings.NewEmbedder(ollamaClient)
	if err != nil {
		panic(err)
	}

	wvStore, err := weaviate.New(
		weaviate.WithEmbedder(emb),
		weaviate.WithScheme(cfg.Scheme),
		weaviate.WithHost(fmt.Sprintf("%s:%s", cfg.WvHost, cfg.WvPort)),
		weaviate.WithIndexName(cfg.IndexName),
	)
	if err != nil {
		panic(err)
	}

	s := &server.RagServer{
		Ctx:          ctx,
		Config:       cfg,
		WvStore:      wvStore,
		Router:       gin.Default(),
		OllamaClient: ollamaClient,
	}
	s.Run()
}

package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
	"log"
	"net/http"
	"strings"
)

type RagServer struct {
	Ctx          context.Context
	Config       *config.Config
	Router       *gin.Engine
	WvStore      weaviate.Store
	OllamaClient *ollama.LLM
}

func (rs *RagServer) Run() error {
	rs.registerRoutes()

	return rs.Router.Run(fmt.Sprintf(":%s", rs.Config.AppPort))
}

func (rs *RagServer) addDocuments(c *gin.Context) {
	var rb addDocumentsRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wvDocs := make([]schema.Document, len(rb.Documents))
	for i, doc := range rb.Documents {
		wvDocs[i] = schema.Document{
			PageContent: doc.Text,
		}
	}

	_, err = rs.WvStore.AddDocuments(rs.Ctx, wvDocs)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error adding documents", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "detail": "documents added successfully"})
}

func (rs *RagServer) query(c *gin.Context) {
	var rb queryRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docs, err := rs.WvStore.SimilaritySearch(rs.Ctx, rb.Content, rs.Config.DocumentsRetrievalNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error querying documents", "error": err.Error()})
		return
	}
	docContents := make([]string, len(docs))
	for i, doc := range docs {
		docContents[i] = doc.PageContent
	}

	ragQuery := fmt.Sprintf(ragTemplateStr, rb.Content, strings.Join(docContents, "\n"))
	respText, err := llms.GenerateFromSinglePrompt(rs.Ctx, rs.OllamaClient, ragQuery, llms.WithModel(rs.Config.LLM))
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error querying documents", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": respText})
}

func (rs *RagServer) registerRoutes() {
	rs.Router.POST("/rag/add", rs.addDocuments)
	rs.Router.POST("/rag/query", rs.query)
}

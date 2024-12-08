package server

import (
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/models"
	"log"
	"net/http"
)

func (s *Server) registerRoutes() {
	s.Router.POST("/rag/add", s.addDocumentsRoute)
	s.Router.POST("/rag/query", s.queryRoute)
}

func (s *Server) addDocumentsRoute(c *gin.Context) {
	var rb models.AddDocumentsRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Rag.AddDocuments(rb)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error adding documents", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "detail": "documents added successfully"})
}

func (s *Server) queryRoute(c *gin.Context) {
	var rb models.QueryRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go s.Rag.Query(rb)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "detail": "query is processed"})
}

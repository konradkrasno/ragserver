package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/konradkrasno/ragserver/models"
	"io"
	"log"
	"net/http"
)

func (s *Server) registerRoutes() {
	s.Router.POST("/rag/api/add", s.addDocumentsRoute)
	s.Router.POST("/rag/api/addFromUrl", s.addDocumentsFromUrlsRoute)
	s.Router.POST("/rag/api/query", s.queryRoute)
	s.Router.GET("/rag/ws/:sessionId", s.wsSendBrokerMessages)
}

func (s *Server) addDocumentsRoute(c *gin.Context) {
	var rb models.AddDocumentsRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go s.Rag.AddDocuments(rb)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "detail": "documents scheduled to be added"})
}

func (s *Server) addDocumentsFromUrlsRoute(c *gin.Context) {
	var rb models.AddDocumentsFromUrlsRequest
	err := c.Bind(&rb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adr := models.AddDocumentsRequest{
		Documents: make([]models.Document, len(rb.Urls)),
	}
	for i, url := range rb.Urls {
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error fetching url", "error": err.Error()})
			return
		}
		body, err := io.ReadAll(resp.Body)

		adr.Documents[i] = models.Document{
			Text: string(body),
		}
	}

	go s.Rag.AddDocuments(adr)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "detail": "documents scheduled to be added"})
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

func (s *Server) wsSendBrokerMessages(c *gin.Context) {
	sessionId := c.Param("sessionId")
	if sessionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
	}

	conn, err := s.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	err = s.Broker.Listen(s.Envs.RabbitMQAnswerExchange, sessionId, func(msg []byte) error {
		return conn.WriteMessage(websocket.TextMessage, msg)
	})
	if err != nil {
		log.Println(err)
		return
	}
}

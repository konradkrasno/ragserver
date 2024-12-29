package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/konradkrasno/ragserver/models"
	"log"
	"net/http"
)

func (s *Server) registerRoutes() {
	s.Router.POST("/api/rag/add", s.addDocumentsRoute)
	s.Router.POST("/api/rag/query", s.queryRoute)
	s.Router.GET("/ws/:sessionId", s.wsSendBrokerMessages)
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

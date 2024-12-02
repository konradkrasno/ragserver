package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/konradkrasno/ragserver/broker"
	"github.com/konradkrasno/ragserver/config"
	"github.com/konradkrasno/ragserver/rag"
)

type Server struct {
	Config   *config.Config
	Router   *gin.Engine
	Rag      *rag.Rag
	Upgrader *websocket.Upgrader
	Broker   broker.Broker
}

func New(cfg *config.Config, router *gin.Engine, rag *rag.Rag, broker broker.Broker) *Server {
	return &Server{
		Config: cfg,
		Router: router,
		Rag:    rag,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Broker: broker,
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.Router.Run(fmt.Sprintf(":%s", s.Config.AppPort))
}

package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/rag"
)

type Server struct {
	Router *gin.Engine
	Rag    *rag.Rag
}

func New(router *gin.Engine, rag *rag.Rag) *Server {
	return &Server{
		Router: router,
		Rag:    rag,
	}
}

func (s *Server) Run(appPort string) error {
	s.registerRoutes()
	return s.Router.Run(fmt.Sprintf(":%s", appPort))
}

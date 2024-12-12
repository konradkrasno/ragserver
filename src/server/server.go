package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/konradkrasno/ragserver/broker"
	"github.com/konradkrasno/ragserver/environment"
	"github.com/konradkrasno/ragserver/rag"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Envs     *environment.Envs
	Router   *gin.Engine
	Rag      *rag.Rag
	Upgrader *websocket.Upgrader
	Broker   broker.Broker
}

func New(envs *environment.Envs, rag *rag.Rag, broker broker.Broker) *Server {
	return &Server{
		Envs:     envs,
		Router:   getRouter(),
		Rag:      rag,
		Upgrader: getUpgrader(envs),
		Broker:   broker,
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.Router.Run(fmt.Sprintf(":%s", s.Envs.AppPort))
}

func getRouter() *gin.Engine {
	router := gin.Default()
	configureRouter(router)

	return router
}

func configureRouter(router *gin.Engine) {
	routerConfig := cors.DefaultConfig()

	routerConfig.AllowAllOrigins = true
	routerConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	routerConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	routerConfig.ExposeHeaders = []string{"Content-Length"}
	routerConfig.AllowCredentials = true
	routerConfig.MaxAge = 12 * time.Hour

	router.Use(cors.New(routerConfig))
}

func getUpgrader(envs *environment.Envs) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if envs.AllowedOrigins == "*" {
				return true
			}

			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range strings.Split(envs.AllowedOrigins, ",") {
				if allowedOrigin == origin {
					return true
				}
			}

			return false
		},
	}
}

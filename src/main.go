package main

import (
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/rag"
	"github.com/konradkrasno/ragserver/server"
	"os"
)

func runServer() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		panic("APP_PORT environment variable not set")
	}

	r, err := rag.New()
	if err != nil {
		panic(err)
	}
	s := server.New(gin.Default(), r)

	err = s.Run(appPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}

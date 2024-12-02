package main

import (
	"github.com/gin-gonic/gin"
	"github.com/konradkrasno/ragserver/broker"
	"github.com/konradkrasno/ragserver/config"
	"github.com/konradkrasno/ragserver/rag"
	"github.com/konradkrasno/ragserver/server"
)

func runServer() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	b := broker.NewMQBroker(cfg.MQEndpoint)

	r, err := rag.New(cfg, b)
	if err != nil {
		panic(err)
	}

	s := server.New(cfg, gin.Default(), r, b)

	err = s.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}

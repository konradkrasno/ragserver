package main

import (
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

	s := server.New(cfg, r, b)

	err = s.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}

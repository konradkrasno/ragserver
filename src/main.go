package main

import (
	"github.com/konradkrasno/ragserver/broker"
	"github.com/konradkrasno/ragserver/environment"
	"github.com/konradkrasno/ragserver/rag"
	"github.com/konradkrasno/ragserver/server"
)

func runServer() {
	envs := environment.NewEnvs()

	b := broker.NewMQBroker(envs)

	r, err := rag.New(envs, b)
	if err != nil {
		panic(err)
	}

	s := server.New(envs, r, b)

	err = s.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer()
}

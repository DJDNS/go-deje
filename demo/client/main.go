package main

import (
	"github.com/campadrenalin/go-deje"
	"log"
	"os"
	"os/signal"
)

func main() {
	url := "ws://localhost:8080/ws"
	topic := "http://demo/"
	logger := log.New(os.Stderr, "deje.SimpleClient: ", 0)

	sc := deje.NewSimpleClient(topic, logger)
	if err := sc.Connect(url); err != nil {
		logger.Fatal(err)
	} else {
		logger.Println("Connected to " + url)
		logger.Printf("Listening to topic '%s'", topic)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down client")
}

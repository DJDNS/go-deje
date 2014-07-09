package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/campadrenalin/go-deje"
)

var host = flag.String("host", "localhost:8080", "Router to connect to")
var topic = flag.String("topic", "http://demo/", "DEJE topic to subscribe to")

func main() {
	flag.Parse()
	url := "ws://" + *host + "/ws"
	logger := log.New(os.Stderr, "deje.SimpleClient: ", 0)

	sc := deje.NewSimpleClient(*topic, logger)
	if err := sc.Connect(url); err != nil {
		logger.Fatal(err)
	} else {
		logger.Println("Connected to " + url)
		logger.Printf("Listening to topic '%s'", *topic)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down client")
}

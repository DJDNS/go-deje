package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/DJDNS/go-deje"
	state "github.com/DJDNS/go-deje/state"
)

var host = flag.String("host", "localhost:8080", "Router to connect to")
var topic = flag.String("topic", "deje://demo/", "DEJE topic to subscribe to")
var filename = flag.String("file", "", "File to load from and save to")

func load(sc *deje.SimpleClient) {
	if *filename == "" {
		return
	}
	log.Printf("Loading from '%s'...", *filename)

	file, err := os.Open(*filename)
	if err != nil {
		log.Printf("Could not load from '%s' - not fatal, though.", *filename)
		return
	}
	defer file.Close()

	doc := sc.GetDoc()
	if err = doc.Deserialize(file); err != nil {
		log.Fatal(err)
	}
	doc.Topic = *topic
	log.Printf("Topic: %s", doc.Topic)
}

func save(sc *deje.SimpleClient) {
	if *filename == "" {
		return
	}
	log.Printf("Writing to '%s'...", *filename)

	file, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc := sc.GetDoc()
	if err = doc.Serialize(file); err != nil {
		log.Fatal(err)
	}
}

type Closer chan struct{}

// Uses booleans for "saturation", coalescing multiple identical requests
func io_ratelimit_loop(sc *deje.SimpleClient, command chan string, closer Closer) {
	var need_load, need_save, need_close bool
	timer := time.Tick(time.Second)

	for {
		select {
		case <-timer:
			if need_load {
				load(sc)
			}
			if need_save {
				save(sc)
			}
			if need_close {
				log.Println("IO rate limiter drained")
				close(closer)
				return
			}
			need_load = false
			need_save = false
		case c := <-command:
			if c == "close" {
				need_close = true
			} else if c == "load" {
				need_load = true
			} else if c == "save" {
				need_save = true
			}
		}
	}
}

func main() {
	flag.Parse()
	url := "ws://" + *host + "/ws"
	logger := log.New(os.Stderr, "deje.SimpleClient: ", 0)
	io_loop_commander := make(chan string)
	io_loop_closer := make(Closer)

	sc := deje.NewSimpleClient(*topic, logger)
	if err := sc.Connect(url); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to " + url)
		log.Printf("Listening to topic '%s'", *topic)
	}

	go io_ratelimit_loop(sc, io_loop_commander, io_loop_closer)
	io_loop_commander <- "load"
	sc.SetPrimitiveCallback(func(p state.Primitive) {
		io_loop_commander <- "save"
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down client")

	io_loop_commander <- "close"
	<-io_loop_closer
}

package main

import (
	"log"
	"os"

	"github.com/DJDNS/go-deje/djconvert/app"
)

func main() {
	logger := log.New(os.Stderr, "djconvert: ", 0)
	if err := app.Main(nil, true); err != nil {
		logger.Fatal(err)
	}
}

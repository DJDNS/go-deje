package main

import (
	"os"

	"github.com/DJDNS/go-deje/djconvert/app"
)

func main() {
	app.Main(nil, true, os.Stderr)
}

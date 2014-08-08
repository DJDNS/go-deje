package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("Insufficient arguments")
	}
	fmt.Printf("Output")
}

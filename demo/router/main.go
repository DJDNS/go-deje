package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/jcelliott/turnpike"
)

func main() {
	server := turnpike.NewServer()
	gopath_dir := os.Getenv("GOPATH")
	host_location := []string{
		gopath_dir,
		"src", "github.com", "DJDNS", "go-deje",
		"demo", "browser",
	}
	static_path := path.Join(host_location...)

	http.Handle("/ws", server.Handler)
	http.Handle("/", http.FileServer(http.Dir(static_path)))

	log.Println("Listening on port 8080")
	log.Println("Serving files from " + static_path)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

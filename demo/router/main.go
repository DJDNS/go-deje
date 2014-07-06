package main

import (
	"bitbucket.org/kardianos/osext"
	"github.com/jcelliott/turnpike"
	"log"
	"net/http"
	"path"
)

func main() {
	server := turnpike.NewServer()
	exc_path, _ := osext.Executable()
	static_path := path.Join(path.Dir(exc_path), "..", "browser")

	http.Handle("/ws", server.Handler)
	http.Handle("/", http.FileServer(http.Dir(static_path)))

	log.Println("Listening on port 8080")
	log.Println("Serving files from " + static_path)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

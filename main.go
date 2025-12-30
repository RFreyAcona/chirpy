package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

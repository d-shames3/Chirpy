package main

import (
	"log"
	"net/http"
)

func main() {
	const filePath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePath)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}

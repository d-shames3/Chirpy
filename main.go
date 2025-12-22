package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = mux
	log.Fatal(server.ListenAndServe())

}

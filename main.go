package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	const healthOKResponse = "OK"

	cfg := &apiConfig{}
	mux := http.NewServeMux()

	healthHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(healthOKResponse))
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("/healthz", healthHandlerFunc)
	mux.HandleFunc("/metrics", cfg.fileServerMetricsHandler)
	mux.HandleFunc("/reset", cfg.resetMetricsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) fileServerMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
}

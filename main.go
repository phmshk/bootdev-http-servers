package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiConfig := &apiConfig{}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))
	prefixHandler := http.StripPrefix("/app/", fs)

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(prefixHandler))
	mux.HandleFunc("GET /api/healthz", firstHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.requestsMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetMetricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	log.Println("Server starting on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
}

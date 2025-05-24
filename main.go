package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) GetCount(w http.ResponseWriter, r *http.Request) {
	hitCount := cfg.fileserverHits.Load()
	// write it to the response in format "Hits: x"

	_, err := w.Write([]byte(fmt.Sprintf("Hits: %d", hitCount)))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) ResetCount(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	_, err := w.Write([]byte(fmt.Sprintf("Reset hit count")))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	//Convert the current directory to a directory for the FileServer
	fileServer := http.FileServer(http.Dir("."))

	api := apiConfig{}

	// Add the file server as a handler for the root path
	mux.Handle("/app/", api.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))

	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/metrics", api.GetCount)
	mux.HandleFunc("/reset", api.ResetCount)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the status code
	w.WriteHeader(http.StatusOK)

	// Write the response body
	_, err := w.Write([]byte("OK"))
	if err != nil {
		// Handle error if writing fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

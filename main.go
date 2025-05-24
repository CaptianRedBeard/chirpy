package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//Convert the current directory to a directory for the FileServer
	fileServer := http.FileServer(http.Dir("."))

	// Add the file server as a handler for the root path
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}

}

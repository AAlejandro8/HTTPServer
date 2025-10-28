package main

import (
	"fmt"
	"log"
	"net/http"
)

func muxHandler(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", muxHandler)

	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("error: %w", err)
	}

}
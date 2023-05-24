package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"goravoxy/github.com/gorilla/mux"
)

func main() {
	// Create a new router using Gorilla Mux
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Define the mappings for the reverse proxy
	mappings := map[string]string{
		"localhost:8080":       "http://localhost:9000",
		"127.0.0.1:8080":       "http://localhost:1986",
		"3dogsfarm.co.za:80":   "http://localhost:9000",
		"boldgear.capetown:80": "http://localhost:1986",
	}

	// Create the reverse proxy handler for each mapping
	for incoming, outgoing := range mappings {
		// Parse the target URL
		targetURL, err := url.Parse(outgoing)
		if err != nil {
			log.Fatalf("Failed to parse target URL for %s: %v", incoming, err)
		}

		// Create the reverse proxy with the target URL
		fmt.Printf("Mapping %s => %s\n", incoming, outgoing)
		fmt.Printf("Target URL: %s\n", targetURL)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Register the reverse proxy handler with the router
		router.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
			// Update the request host to match the target host
			r.Host = targetURL.Host

			// Serve the request through the reverse proxy
			proxy.ServeHTTP(w, r)
		}).Host(incoming)
	}

	// Start the HTTP server
	log.Println("Reverse proxy is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

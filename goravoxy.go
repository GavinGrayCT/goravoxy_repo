package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"goravoxy/github.com/gorilla/mux"
)

func main() {
	var hostA *string = flag.String("a", "localhost", "Host a url.")
	var hostB *string = flag.String("b", "localhost", "Host b url.")
	var portServing *string = flag.String("p", "80", "Serving Port Address.")
	var targetAPort *string = flag.String("t", "9000", "Target port for Host a.")
	var targetBPort *string = flag.String("u", "9001", "Target port for Host b.")
	flag.Parse()

	// Create a new router using Gorilla Mux
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Define the mappings for the reverse proxy
	mappings := map[string]string{
		fmt.Sprintf("%s:%s", *hostA, *portServing): fmt.Sprintf("http://localhost:%s", *targetAPort),
		fmt.Sprintf("%s:%s", *hostB, *portServing): fmt.Sprintf("http://localhost:%s", *targetBPort),
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

		proxy.ServeHTTP(w, r)
	}

	// Start the HTTP server
	log.Printf("Reverse proxy is running on :%s", *portServing)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *portServing), router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Deliberate flaw for SAST tool to find later
const hardcodedApiKey = "dummy-go-key-9876"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Using the hardcoded key just to show SAST detection
	log.Printf("Request received. Using key (don't log keys!): %s", hardcodedApiKey)
	fmt.Fprintf(w, "Hello, Secure Go World!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	http.HandleFunc("/", helloHandler)

	log.Printf("Server starting on port %s\n", port)
	// Use the address format required by http.ListenAndServe
	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

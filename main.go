package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Deliberate flaw for SAST tool to find later
// const hardcodedApiKey = "dummy-go-key-9876"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Load the API Key from environment variable within the handler
	// Or load it once in main() and pass it down if used frequently
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("Warning: API_KEY environment variable not set.")
		apiKey = "not-set" // Provide a default or handle error appropriately
	}

	// Using the hardcoded key just to show SAST detection
	log.Printf("Request received. Using key (don't log keys!): %s", apiKey)
	fmt.Fprintf(w, "Hello, Secure Go World!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Check for TLS configuration via environment variables
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	http.HandleFunc("/", helloHandler)

	log.Printf("Server starting on port %s\n", port)
	// Use the address format required by http.ListenAndServe
	addr := fmt.Sprintf(":%s", port)
	// if err := http.ListenAndServe(addr, nil); err != nil {
	// 	log.Fatalf("Error starting server: %s\n", err)
	// }

	if certFile != "" && keyFile != "" {
		// If TLS cert and key file paths are provided, start HTTPS server
		log.Printf("TLS cert and key files provided. Starting HTTPS server on port %s\n", port)
		// Note: Ensure the certFile and keyFile paths are accessible within the container/environment
		if err := http.ListenAndServeTLS(addr, certFile, keyFile, nil); err != nil {
			log.Fatalf("Error starting HTTPS server: %s\n", err)
		}
	} else {
		// If TLS files are not provided, start standard HTTP server
		log.Printf("TLS cert and key files not found in env vars. Starting HTTP server on port %s\n", port)
		// This might still be flagged by Semgrep if the rule is strict,
		// but it addresses the finding by attempting TLS first.
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Error starting HTTP server: %s\n", err)
		}
	}
}

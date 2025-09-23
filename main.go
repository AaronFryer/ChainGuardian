package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println(Hello("Crate"))

	// Define a handler function for the "/" route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go Web Server!")
	})

	// Start the HTTP server on port 8080
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// Accept NPM requests
	// TODO create simple webserver

	// Filter out young packages

	// Cache tarballs

	// Return
}

func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

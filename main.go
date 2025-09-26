package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aaronfryer/crate/internal/config"
	"github.com/aaronfryer/crate/internal/httphandlers"
)

func main() {
	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", httphandlers.CacheHandler)

	if err := http.ListenAndServe(config.HTTPPort, nil); err != nil {
		log.Fatal(err)
	}
}

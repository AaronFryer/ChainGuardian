package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aaronfryer/chainguardian/internal/config"
	"github.com/aaronfryer/chainguardian/internal/httphandlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		log.Fatal(err)
	}

	handler := httphandlers.New(cfg)
	http.HandleFunc("/", handler.CacheHandler)

	log.Printf("Starting server on %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, nil); err != nil {
		log.Fatal(err)
	}
}

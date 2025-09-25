package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	cacheDir       = "./cache"
	remoteRegistry = "https://registry.npmjs.org"
	httpPort       = ":8080"
)

// TODO
// Additional registries
// logging & metrics
// health check endpoint
// test everything
// ban young package versions

func HandlePackageJsonRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(remoteRegistry + r.URL.Path)

	if err != nil {
		http.Error(w, "Failed to fetch package info", http.StatusBadGateway)
		log.Printf("Package fetch error: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Upstream server error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		log.Printf("Body read error: %v", err)
		return
	}

	packageJsonFilePath := cacheDir + r.URL.Path + ".json"
	if err := os.WriteFile(packageJsonFilePath, body, os.FileMode(0644)); err != nil {
		log.Printf("Cache package write error: %v", err)
	}

	// TODO serve from cache only

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func HandleTarballRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(remoteRegistry + r.URL.Path)

	if err != nil {
		http.Error(w, "Failed to fetch tarball", http.StatusBadGateway)
		log.Printf("Tarball fetch error: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Upstream server error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		log.Printf("Body read error: %v", err)
		return
	}

	fileName := strings.Split(r.URL.Path, "/")

	if err := os.MkdirAll(cacheDir+"/"+fileName[1], 0755); err != nil {
		log.Fatal(err)
	}
	filePath := cacheDir + "/" + fileName[1] + "/" + fileName[3]

	if err := os.WriteFile(filePath, body, os.FileMode(0644)); err != nil {
		log.Printf("Cache tarball write error: %v", err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s url: %s", r.Method, r.URL.Path)

	if CountPathSegments(r.URL.Path) == 1 {
		HandlePackageJsonRequest(w, r)

	} else if CountPathSegments(r.URL.Path) == 3 && strings.HasSuffix(r.URL.Path, ".tgz") {
		HandleTarballRequest(w, r)

	} else {
		fmt.Fprintln(w, "Unsupported request")
	}

}

func main() {
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server
	http.HandleFunc("/", cacheHandler)

	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
// todo handle timeouts

func isVersionOldEnough(timestamp string) bool {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false // If we can't parse the time, reject it
	}
	age := time.Since(t)
	return age >= 60*24*time.Hour && age <= 365*24*time.Hour
}

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

	// TODO Do filtering
	var packageData map[string]interface{}
	if err := json.Unmarshal(body, &packageData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		log.Printf("JSON parse error: %v", err)
		return
	}

	if timeData, ok := packageData["time"].(map[string]interface{}); ok {
		filteredTime := make(map[string]interface{})
		filteredVersions := make(map[string]interface{})

		// TODO somehow log when someone tries to access an old or young version
		for version, timestamp := range timeData {
			if ts, ok := timestamp.(string); ok && isVersionOldEnough(ts) {
				filteredTime[version] = timestamp
				if versions, ok := packageData["versions"].(map[string]interface{}); ok {
					if versionData, exists := versions[version]; exists {
						filteredVersions[version] = versionData
					}
				}
			}
		}

		packageData["time"] = filteredTime
		packageData["versions"] = filteredVersions
	}

	modifiedBody, err := json.Marshal(packageData)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
		return
	}

	packageJsonFilePath := cacheDir + r.URL.Path + ".json"
	if err := os.WriteFile(packageJsonFilePath, modifiedBody, os.FileMode(0644)); err != nil {
		log.Printf("Cache package write error: %v", err)
	}

	// TODO serve from cache only

	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, packageJsonFilePath)

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

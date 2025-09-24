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
)

func countPathSegments(path string) int {
	trimmedPath := strings.Trim(path, "/")
	if trimmedPath == "" {
		return 0
	}
	return len(strings.Split(trimmedPath, "/"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s url: %s", r.Method, r.URL.Path)

	if countPathSegments(r.URL.Path) == 1 {
		resp, _ := http.Get(remoteRegistry + r.URL.Path)

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		os.WriteFile(cacheDir+r.URL.Path+".json", body, os.FileMode(0644))

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else if countPathSegments(r.URL.Path) == 3 && strings.HasSuffix(r.URL.Path, ".tgz") {
		resp, _ := http.Get(remoteRegistry + r.URL.Path)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fileName := strings.Split(r.URL.Path, "/")

		if err := os.MkdirAll(cacheDir+"/"+fileName[1], 0755); err != nil {
			log.Fatal(err)
		}
		filePath := cacheDir + "/" + fileName[1] + "/" + fileName[3]
		os.WriteFile(filePath, body, os.FileMode(0644))

		// serve the file from the os cache folder
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filePath)

	} else {
		fmt.Fprintln(w, "Hello, World!")
	}

}

func main() {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

package httphandlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aaronfryer/crate/internal/cache"
	"github.com/aaronfryer/crate/internal/config"
)

type HttpHandler struct {
	cfg *config.Config
}

// New creates a new Handler instance with the provided configuration
func New(cfg *config.Config) *HttpHandler {
	return &HttpHandler{
		cfg: cfg,
	}
}

func (h *HttpHandler) HandlePackageJSON(w http.ResponseWriter, r *http.Request) {
	RemoteUrl := h.cfg.RemoteRegistry

	if registryUrl, ok := h.cfg.Registries[strings.TrimPrefix(r.URL.Path, "/")]; ok {
		RemoteUrl = registryUrl
	}

	resp, err := http.Get(RemoteUrl + r.URL.Path)

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

	var packageData map[string]interface{}
	if err := json.Unmarshal(body, &packageData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		log.Printf("JSON parse error: %v", err)
		return
	}

	h.filterPackageVersions(packageData)

	modifiedBody, err := json.Marshal(packageData)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		log.Printf("JSON encode error: %v", err)
		return
	}

	packageJsonFilePath := h.cfg.CacheDir + r.URL.Path + ".json"
	cache.SavePackageJSON(h.cfg, r.URL.Path, modifiedBody)

	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, packageJsonFilePath)
}

func (h *HttpHandler) HandleTarball(w http.ResponseWriter, r *http.Request) {
	RemoteUrl := h.cfg.RemoteRegistry

	if registryUrl, ok := h.cfg.Registries[strings.TrimPrefix(r.URL.Path, "/")]; ok {
		RemoteUrl = registryUrl
	}

	resp, err := http.Get(RemoteUrl + r.URL.Path)

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
	cache.SaveTarball(h.cfg, fileName[1], fileName[3], body)

	w.Header().Set("Content-Type", "application/octet-stream")
	filePath := h.cfg.CacheDir + "/" + fileName[1] + "/" + fileName[3]
	http.ServeFile(w, r, filePath)
}

func (h *HttpHandler) CacheHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s url: %s", r.Method, r.URL.Path)

	segments := strings.Count(r.URL.Path, "/")

	switch {
	case segments == 1:
		h.HandlePackageJSON(w, r)
	case segments == 3 && strings.HasSuffix(r.URL.Path, ".tgz"):
		h.HandleTarball(w, r)
	default:
		http.Error(w, "Unsupported request", http.StatusBadRequest)
	}
}

func (h *HttpHandler) isVersionOldEnough(timestamp string) bool {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false
	}
	age := time.Since(t)
	return age >= h.cfg.MinPackageAge && age <= h.cfg.MaxPackageAge
}

func (h *HttpHandler) filterPackageVersions(packageData map[string]interface{}) {
	timeData, ok := packageData["time"].(map[string]interface{})
	if !ok {
		return
	}

	filteredTime := make(map[string]interface{})
	filteredVersions := make(map[string]interface{})

	for version, timestamp := range timeData {
		if ts, ok := timestamp.(string); ok && h.isVersionOldEnough(ts) {
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

package cache

import (
	"os"
	"path/filepath"

	"github.com/aaronfryer/crate/internal/config"
)

func SavePackageJSON(cfg *config.Config, path string, data []byte) error {
	fullPath := filepath.Join(cfg.CacheDir, path+".json")
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0644)
}

func SaveTarball(cfg *config.Config, packageName string, fileName string, data []byte) error {
	cleanPackageName := filepath.Base(packageName)
	cleanFileName := filepath.Base(fileName)

	dir := filepath.Join(cfg.CacheDir, cleanPackageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, cleanFileName), data, 0644)
}

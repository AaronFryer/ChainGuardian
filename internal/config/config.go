package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	CacheDir       string        `toml:"cache_dir"`
	RemoteRegistry string        `toml:"remote_registry"`
	HTTPPort       string        `toml:"http_port"`
	MinPackageAge  time.Duration `toml:"-"`
	MaxPackageAge  time.Duration `toml:"-"`

	// Raw values from TOML
	MinPackageAgeHours int64 `toml:"min_package_age"`
	MaxPackageAgeHours int64 `toml:"max_package_age"`
}

var DefaultConfig = Config{
	CacheDir:       "./cache",
	RemoteRegistry: "https://registry.npmjs.org",
	HTTPPort:       ":8080",
	MinPackageAge:  60 * 24 * time.Hour,
	MaxPackageAge:  15 * 365 * 24 * time.Hour,
}

// Load reads the configuration from a TOML file
func Load() (*Config, error) {
	// Get the directory of the executable
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	configPath := filepath.Join(projectRoot, "config.toml")

	config := DefaultConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Printf("Warning: Could not load config file, using defaults: %v", err)
		return &DefaultConfig, nil
	}

	// Convert hours to duration
	config.MinPackageAge = time.Duration(config.MinPackageAgeHours) * time.Hour
	config.MaxPackageAge = time.Duration(config.MaxPackageAgeHours) * time.Hour

	return &config, nil
}

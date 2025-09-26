package config

import "time"

const (
	CacheDir       = "./cache"
	RemoteRegistry = "https://registry.npmjs.org"
	HTTPPort       = ":8080"
	MinPackageAge  = 60 * 24 * time.Hour
	MaxPackageAge  = 15 * 365 * 24 * time.Hour
)

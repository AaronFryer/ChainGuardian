package utils

import "strings"

func CountPathSegments(path string) int {
	trimmedPath := strings.Trim(path, "/")
	if trimmedPath == "" {
		return 0
	}
	return len(strings.Split(trimmedPath, "/"))
}

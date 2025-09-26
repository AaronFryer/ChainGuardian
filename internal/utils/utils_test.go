package utils

import "testing"

func TestCountPathSegments(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "empty path",
			path:     "",
			expected: 0,
		},
		{
			name:     "root path",
			path:     "/",
			expected: 0,
		},
		{
			name:     "single segment",
			path:     "/express",
			expected: 1,
		},
		{
			name:     "single segment without slashes",
			path:     "express",
			expected: 1,
		},
		{
			name:     "multiple segments",
			path:     "/express/express-1.0.0.tgz",
			expected: 2,
		},
		{
			name:     "multiple segments with trailing slash",
			path:     "/express/latest/",
			expected: 2,
		},
		{
			name:     "package download path",
			path:     "/express/express-1.0.0/express-1.0.0.tgz",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountPathSegments(tt.path)
			if result != tt.expected {
				t.Errorf("CountPathSegments(%q) = %d, want %d", tt.path, result, tt.expected)
			}
		})
	}
}

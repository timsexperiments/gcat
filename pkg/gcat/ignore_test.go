package gcat

import (
	"path/filepath"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	patterns := []string{
		"*.log",
		"temp/",
		"vendor/",
	}

	tests := []struct {
		relPath string
		want    bool
	}{
		{
			relPath: "error.log",
			want:    true,
		},
		{
			relPath: "temp/file.txt",
			want:    true,
		},
		{
			relPath: "vendor/package/file.go",
			want:    true,
		},
		{
			relPath: "src/main.go",
			want:    false,
		},
		{
			relPath: "doc/readme.md",
			want:    false,
		},
		{
			relPath: "tempdir/file.txt",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			got := shouldIgnore(filepath.ToSlash(tt.relPath), patterns)
			if got != tt.want {
				t.Errorf("shouldIgnore(%q, %v) = %v; want %v", tt.relPath, patterns, got, tt.want)
			}
		})
	}
}

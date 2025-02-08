package gcat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoCommon_GetLanguage(t *testing.T) {
	t.Parallel()

	rc := newRepoCommon()

	tests := []struct {
		filename string
		want     string
	}{
		// Programming Languages
		{
			filename: "main.go",
			want:     "Go",
		},
		{
			filename: "script.py",
			want:     "Python",
		},
		{
			filename: "index.js",
			want:     "JavaScript",
		},
		{
			filename: "component.jsx",
			want:     "JavaScript (React)",
		},
		{
			filename: "module.ts",
			want:     "TypeScript",
		},
		{
			filename: "module.tsx",
			want:     "TypeScript (React)",
		},
		{
			filename: "Program.java",
			want:     "Java",
		},
		{
			filename: "hello.c",
			want:     "C",
		},
		{
			filename: "hello.cpp",
			want:     "C++",
		},
		{
			filename: "app.cs",
			want:     "C#",
		},
		{
			filename: "ruby.rb",
			want:     "Ruby",
		},
		{
			filename: "server.rs",
			want:     "Rust",
		},
		// Markup / Web
		{
			filename: "index.html",
			want:     "HTML",
		},
		{
			filename: "style.css",
			want:     "CSS",
		},
		{
			filename: "readme.md",
			want:     "Markdown",
		},
		// Configuration & Data Formats
		{
			filename: "config.json",
			want:     "JSON",
		},
		{
			filename: "settings.yaml",
			want:     "YAML",
		},
		{
			filename: "data.toml",
			want:     "TOML",
		},
		// Scripting Languages
		{
			filename: "script.sh",
			want:     "Shell Script",
		},
		// Files without a known extension should return an empty string.
		{
			filename: "unknown.xyz",
			want:     "",
		},
		// Case-insensitive test.
		{
			filename: "README.MD",
			want:     "Markdown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()

			got := rc.getLanguage(tt.filename)
			if got != tt.want {
				t.Errorf("getLanguage(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestWithRegisteredLanguages(t *testing.T) {
	t.Parallel()

	overrides := map[string]string{
		".go":   "Golang",
		".test": "Test Script",
	}

	WithRegisteredLanguages(overrides)

	tests := []struct {
		ext  string
		want string
	}{
		{
			ext: ".go", want: "Golang",
		},
		{
			ext: ".test", want: "Test Script",
		},
		{
			ext: ".py", want: "Python",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("for file %s", tt.ext), func(t *testing.T) {
			t.Parallel()

			repo, err := OpenRepository("testdata/no-ignore", WithRegisteredLanguages(overrides))
			require.NoError(t, err)
			filename := "file" + tt.ext
			got := repo.GetLanguage(filename)
			if got != tt.want {
				t.Errorf("After override, getLanguage(%q) = %q, want %q", filename, got, tt.want)
			}
		})
	}
}

func TestNewRepoCommon(t *testing.T) {
	rc := newRepoCommon()
	if rc.languages == nil {
		t.Fatalf("Expected languages map to be non-nil")
	}

	defaults := []struct {
		ext  string
		want string
	}{
		{
			ext:  ".go",
			want: "Go",
		},
		{
			ext:  ".py",
			want: "Python",
		},
		{
			ext:  ".html",
			want: "HTML",
		},
		{
			ext:  ".md",
			want: "Markdown",
		},
	}

	for _, d := range defaults {
		t.Run(fmt.Sprintf("for %s", d.ext), func(t *testing.T) {
			t.Parallel()

			if got, ok := rc.languages[d.ext]; !ok || got != d.want {
				t.Errorf("Default for %q = %q, want %q", d.ext, got, d.want)
			}
		})
	}
}

func TestOpenRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		opts   []Option
		root   string
		assert func(*testing.T, Repository, error)
	}{
		{
			name: "Empty",
			root: "",
			assert: func(t *testing.T, r Repository, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
		{
			name: "Non-existent",
			root: "non-existent",
			assert: func(t *testing.T, r Repository, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
		{
			name: "File",
			root: "testdata/file.txt",
			assert: func(t *testing.T, r Repository, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
		{
			name: "Directory",
			root: "testdata",
			assert: func(t *testing.T, r Repository, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			},
		},
		{
			name: "WithRegisteredLanguages",
			opts: []Option{
				WithRegisteredLanguages(map[string]string{
					".test": "Test",
				}),
			},
			root: "testdata",
			assert: func(t *testing.T, r Repository, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Contains(t, r.GetLanguage("x.test"), "Test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := OpenRepository(tt.root, tt.opts...)

			tt.assert(t, repo, err)
		})
	}
}

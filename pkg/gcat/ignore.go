package gcat

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func loadGitIgnore(dir string) ([]string, error) {
	gitIgnorePath := filepath.Join(dir, ".gitignore")
	file, err := os.Open(gitIgnorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}

func shouldIgnore(relPath string, patterns []string) bool {
	relPath = filepath.ToSlash(relPath)
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(pattern)
		if strings.HasSuffix(pattern, "/") {
			trimmed := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, trimmed+"/") {
				return true
			}
		} else {
			matched, err := doublestar.PathMatch(pattern, relPath)
			if err == nil && matched {
				return true
			}
		}
	}
	return false
}

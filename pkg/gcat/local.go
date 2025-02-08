package gcat

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var defaultIgnore = []string{".git", ".svn", ".hg", ".bzr"}

type localRepository struct {
	root   string
	common *repoCommon
	ignore []string
}

func (l *localRepository) GetFiles() ([]string, error) {
	var files []string
	err := filepath.WalkDir(l.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ignoredContent, err := loadGitIgnore(filepath.Dir(path))
		if err != nil {
			return err
		}
		for _, pattern := range ignoredContent {
			l.registerIgnore(filepath.Join(filepath.Dir(path), pattern))
		}

		if l.shouldIgnore(path) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(l.root, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func (l *localRepository) GetFileContent(filePath string) (string, error) {
	fullPath := filepath.Join(l.root, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (l *localRepository) ConcatFiles(files []string) (string, error) {
	var sb strings.Builder
	sort.Strings(files)
	for i, filePath := range files {
		content, err := l.GetFileContent(filePath)
		if err != nil {
			return "", err
		} else {
			lang := l.common.getLanguage(filePath)
			if lang != "" {
				sb.WriteString(fmt.Sprintf("%s (%s):\n\n", filePath, lang))
			} else {
				sb.WriteString(fmt.Sprintf("%s:\n\n", filePath))
			}
			sb.WriteString("<contents>\n")
			sb.WriteString(content)
			sb.WriteString("\n</contents>")
		}
		if i < len(files)-1 {
			sb.WriteString("\n\n---\n\n")
		}
	}
	return sb.String(), nil
}

func (l *localRepository) GetLanguage(filePath string) string {
	return l.common.getLanguage(filePath)
}

func NewLocalRepository(root string, opts ...Option) (Repository, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}
	repo := &localRepository{root: root, common: newRepoCommon(), ignore: defaultIgnore}
	for _, opt := range opts {
		opt(repo)
	}
	return repo, nil
}

func (l *localRepository) registerIgnore(ignore string) {
	l.ignore = append(l.ignore, ignore)
}

func (l *localRepository) shouldIgnore(filePath string) bool {
	for _, pattern := range l.ignore {
		matched, err := doublestar.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}
	}
	return false
}

package gcat

import (
	"fmt"
	"io"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitRepository struct {
	repo   *git.Repository
	common *repoCommon
}

func (g *gitRepository) GetFiles() ([]string, error) {
	ref, err := g.repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	var files []string
	err = tree.Files().ForEach(func(f *object.File) error {
		files = append(files, f.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (g *gitRepository) GetFileContent(filePath string) (string, error) {
	ref, err := g.repo.Head()
	if err != nil {
		return "", err
	}
	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}
	file, err := tree.File(filePath)
	if err != nil {
		return "", err
	}
	reader, err := file.Blob.Reader()
	if err != nil {
		return "", err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (g *gitRepository) ConcatFiles(files []string) (string, error) {
	var sb strings.Builder
	sort.Strings(files)
	for i, filePath := range files {
		content, err := g.GetFileContent(filePath)
		if err != nil {
			return "", err
		} else {
			lang := g.common.getLanguage(filePath)
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

func (g *gitRepository) GetLanguage(filePath string) string {
	return g.common.getLanguage(filePath)
}

func CloneGitRepository(repoURL string, opts ...Option) (Repository, error) {
	gitRepo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}
	repo := &gitRepository{repo: gitRepo, common: newRepoCommon()}
	for _, opt := range opts {
		opt(repo)
	}
	return repo, nil
}

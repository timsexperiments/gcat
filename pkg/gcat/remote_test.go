package gcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const repository = "https://github.com/timsexperiments/gcat.git"

func TestRemoteRepository_GetFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name: "lists all files from local path",
			want: []string{
				"pkg/gcat/testdata/no-ignore/file.no-language",
				"pkg/gcat/testdata/no-ignore/file1.txt",
				"pkg/gcat/testdata/no-ignore/file2.txt",
				"pkg/gcat/testdata/no-ignore/nested/nested_file.txt",
				"pkg/gcat/testdata/with-ignore/.gitignore",
				"pkg/gcat/testdata/with-ignore/file1.txt",
				"pkg/gcat/testdata/with-ignore/file2.txt",
				"pkg/gcat/testdata/with-ignore/ignored/ignored_file.txt",
				"pkg/gcat/testdata/with-ignore/nested/.gitignore",
				"pkg/gcat/testdata/with-ignore/nested/ignored_file.txt",
				"pkg/gcat/testdata/with-ignore/nested/nested_file.txt",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := OpenRepository(repository)
			require.NoError(t, err)

			files, err := repo.GetFiles()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}
			assert.NoError(t, err)
			for _, file := range tt.want {
				assert.Contains(t, files, file)
			}
		})
	}
}

func TestRemoteRepository_ConcatFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		files   []string
		want    string
		wantErr bool
	}{
		{
			name: "concatenates files from local path",
			path: "testdata/no-ignore",
			files: []string{
				"pkg/gcat/testdata/no-ignore/file1.txt",
				"pkg/gcat/testdata/no-ignore/file2.txt",
				"pkg/gcat/testdata/no-ignore/nested/nested_file.txt",
			},
			want:    "pkg/gcat/testdata/no-ignore/file1.txt (Text):\n\n<contents>\ntest content\n</contents>\n\n---\n\npkg/gcat/testdata/no-ignore/file2.txt (Text):\n\n<contents>\n\n</contents>\n\n---\n\npkg/gcat/testdata/no-ignore/nested/nested_file.txt (Text):\n\n<contents>\n\n</contents>",
			wantErr: false,
		},
		{
			name: "non-existent file",
			path: "testdata/with-ignore",
			files: []string{
				"pkg/gcat/testdata/with-ignore/non-existent/.gitignore",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "includes file with no language",
			files: []string{
				"pkg/gcat/testdata/no-ignore/file.no-language",
			},
			want:    "pkg/gcat/testdata/no-ignore/file.no-language:\n\n<contents>\n\n</contents>",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := OpenRepository(repository)
			require.NoError(t, err)
			files, err := repo.ConcatFiles(tt.files)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, files)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, files)
		})
	}
}

func TestRemoteRepository_GetLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "returns language for file",
			filePath: "file1.txt",
			want:     "Text File",
		},
		{
			name:     "returns empty string for file with no language",
			filePath: "pkg/gcat/testdata/no-ignore/file.no-language",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := OpenRepository(repository, WithRegisteredLanguages(map[string]string{".txt": "Text File"}))
			require.NoError(t, err)

			lang := repo.GetLanguage(tt.filePath)
			assert.Equal(t, tt.want, lang)
		})
	}
}

func TestRemoteRepository_GetFileContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		filePath string
		want     string
		wantErr  bool
	}{
		{
			name:     "returns file content",
			filePath: "pkg/gcat/testdata/no-ignore/file1.txt",
			want:     "test content",
			wantErr:  false,
		},
		{
			name:     "returns error for non-existent file",
			filePath: "pkg/gcat/testdata/no-ignore/non-existent",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := OpenRepository(repository)
			require.NoError(t, err)

			content, err := repo.GetFileContent(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, content)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, content)
		})
	}
}

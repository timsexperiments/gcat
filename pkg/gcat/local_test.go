package gcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalRepository_GetFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "lists all files from local path",
			path: "testdata/no-ignore",
			want: []string{
				"file.no-language",
				"file1.txt",
				"file2.txt",
				"nested/nested_file.txt",
			},
			wantErr: false,
		},
		{
			name: "ignores files from gitignore",
			path: "testdata/with-ignore",
			want: []string{
				".gitignore",
				"nested/.gitignore",
				"nested/nested_file.txt",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := NewLocalRepository(tt.path)
			require.NoError(t, err)

			files, err := repo.GetFiles()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, files)
		})
	}
}

func TestLocalRepository_ConcatFiles(t *testing.T) {
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
				"file1.txt",
				"file2.txt",
				"nested/nested_file.txt",
			},
			want:    "file1.txt (Text):\n\n<contents>\ntest content\n</contents>\n\n---\n\nfile2.txt (Text):\n\n<contents>\n\n</contents>\n\n---\n\nnested/nested_file.txt (Text):\n\n<contents>\n\n</contents>",
			wantErr: false,
		},
		{
			name: "non-existent file",
			path: "testdata/with-ignore",
			files: []string{
				"non-existent/.gitignore",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "includes file with no language",
			path: "testdata/no-ignore",
			files: []string{
				"file.no-language",
			},
			want:    "file.no-language:\n\n<contents>\n\n</contents>",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := NewLocalRepository(tt.path)
			if err != nil {
				t.Fatal(err)
			}
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

func TestLocalRepository_GetLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		opts     []Option
		filePath string
		want     string
	}{
		{
			name:     "returns language for file",
			opts:     []Option{WithRegisteredLanguages(map[string]string{".txt": "Text File"})},
			filePath: "file1.txt",
			want:     "Text File",
		},
		{
			name:     "returns language for file with no language",
			filePath: "file.no-language",
			want:     "",
		},
		{
			name:     "returns language for file with no language in nested directory",
			filePath: "testdata/with-ignore/nested/file.no-language",
			want:     "",
		},
		{
			name:     "returns language for file with no language in nested directory with ignore",
			filePath: "testdata/with-ignore/nested/ignored/ignored_file.txt",
			want:     "Text",
		},
		{
			name:     "returns language for file with no language in nested directory with ignore and nested",
			filePath: "testdata/with-ignore/nested/ignored/nested_file.txt",
			want:     "Text",
		},
		{
			name:     "returns language for file with no language in nested directory with ignore and nested and nested",
			filePath: "testdata/with-ignore/nested/ignored/nested_file.txt",
			want:     "Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := NewLocalRepository("testdata", tt.opts...)
			require.NoError(t, err)

			lang := repo.GetLanguage(tt.filePath)
			assert.Equal(t, tt.want, lang)
		})
	}
}

func TestLocalRepository_GetFileContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		filePath string
		want     string
		wantErr  bool
	}{
		{
			name:     "returns content for file",
			path:     "testdata/no-ignore",
			filePath: "file1.txt",
			want:     "test content",
			wantErr:  false,
		},
		{
			name:     "returns content for file in nested directory",
			path:     "testdata/no-ignore",
			filePath: "non-existent.txt",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := NewLocalRepository(tt.path)
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

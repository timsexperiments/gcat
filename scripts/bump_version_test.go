package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcNewVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		initialTag string
		args       []string
		expected   string
		wantErr    bool
	}{
		{
			name:       "initial patch bump",
			initialTag: "",
			args:       []string{"patch"},
			expected:   "v0.0.1",
			wantErr:    false,
		},
		{
			name:       "initial minor bump",
			initialTag: "",
			args:       []string{"minor"},
			expected:   "v0.1.0",
			wantErr:    false,
		},
		{
			name:       "initial major bump",
			initialTag: "",
			args:       []string{"major"},
			expected:   "v1.0.0",
			wantErr:    false,
		},
		{
			name:       "patch bump normal",
			initialTag: "v1.2.3",
			args:       []string{"patch"},
			expected:   "v1.2.4",
			wantErr:    false,
		},
		{
			name:       "minor bump normal",
			initialTag: "v1.2.3",
			args:       []string{"minor"},
			expected:   "v1.3.0",
			wantErr:    false,
		},
		{
			name:       "major bump normal",
			initialTag: "v1.2.3",
			args:       []string{"major"},
			expected:   "v2.0.0",
			wantErr:    false,
		},
		{
			name:       "patch bump with alpha when base is normal",
			initialTag: "v1.2.3",
			args:       []string{"patch", "alpha"},
			expected:   "v1.2.4-alpha.1",
			wantErr:    false,
		},
		{
			name:       "patch bump with alpha when base is already alpha",
			initialTag: "v1.2.3-alpha.1",
			args:       []string{"patch", "alpha"},
			expected:   "v1.2.3-alpha.2",
			wantErr:    false,
		},
		{
			name:       "patch bump with beta when base is already beta",
			initialTag: "v1.2.3-beta.2",
			args:       []string{"patch", "beta"},
			expected:   "v1.2.3-beta.3",
			wantErr:    false,
		},
		{
			name:       "patch bump with alpha when switching from beta",
			initialTag: "v1.2.3-beta.2",
			args:       []string{"patch", "alpha"},
			expected:   "Error: Switching from beta to alpha in a patch bump is not allowed without a numeric bump.",
			wantErr:    true,
		},
		{
			name:       "patch bump with beta when switching from alpha",
			initialTag: "v1.2.3-alpha.2",
			args:       []string{"patch", "beta"},
			expected:   "v1.2.3-beta.1",
			wantErr:    false,
		},
		{
			name:       "no arguments provided",
			initialTag: "v1.2.3",
			args:       []string{},
			expected:   "Usage: ",
			wantErr:    true,
		},
		{
			name:       "invalid release type",
			initialTag: "v1.2.3",
			args:       []string{"invalid"},
			expected:   "Error: release type must be one of:",
			wantErr:    true,
		},
		{
			name:       "invalid prerelease value",
			initialTag: "v1.2.3",
			args:       []string{"patch", "rc"},
			expected:   "Error: prerelease, if provided, must be alpha or beta",
			wantErr:    true,
		},
		{
			name:       "invalid latest tag format",
			initialTag: "vinvalid",
			args:       []string{"patch"},
			expected:   "Error: Latest tag 'vinvalid' is not in a valid semantic version format.",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tempDir, err := os.MkdirTemp("", "testrepo-*")
			require.NoError(t, err, "failed to create temp dir")
			t.Cleanup(func() { os.RemoveAll(tempDir) })

			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "git init failed: %s", string(out))

			if tt.initialTag != "" {
				cmd = exec.Command("git", "commit", "--allow-empty", "-m", "init")
				cmd.Dir = tempDir
				out, err = cmd.CombinedOutput()
				require.NoError(t, err, "git commit failed: %s", string(out))

				cmd = exec.Command("git", "tag", tt.initialTag)
				cmd.Dir = tempDir
				out, err = cmd.CombinedOutput()
				require.NoError(t, err, "git tag failed: %s", string(out))
			}

			scriptPath, err := filepath.Abs(filepath.Join("bump_version.sh"))
			require.NoError(t, err, "failed to get absolute path for the script")

			cmd = exec.Command(scriptPath, tt.args...)
			cmd.Dir = tempDir
			outputBytes, err := cmd.CombinedOutput()
			fmt.Println(string(outputBytes))

			if tt.wantErr {
				require.Error(t, err, "script unexpectedly succeeded")
				assert.Contains(t, string(outputBytes), tt.expected, "expected error message %s, got %s", tt.expected, string(outputBytes))
			} else {
				require.NoError(t, err, "script failed: %v", err)
				newVersion := strings.TrimSpace(string(outputBytes))
				assert.Contains(t, tt.expected, newVersion, "expected version %s, got %s", tt.expected, newVersion)
			}
		})
	}
}

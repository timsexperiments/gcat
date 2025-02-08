package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleSelector(t *testing.T) {
	realAskOne := askOne
	defer func() { askOne = realAskOne }()

	testCases := []struct {
		name        string
		mockFn      func(p survey.Prompt, response interface{}, options ...survey.AskOpt) error
		inputFiles  []string
		expected    []string
		expectedErr string
	}{
		{
			name: "Successful selection",
			mockFn: func(p survey.Prompt, response interface{}, options ...survey.AskOpt) error {
				// Simulate a user selecting "file1.txt" and "file2.txt"
				if sel, ok := response.(*[]string); ok {
					*sel = []string{"file1.txt", "file2.txt"}
				}
				return nil
			},
			inputFiles: []string{"file1.txt", "file2.txt", "file3.txt"},
			expected:   []string{"file1.txt", "file2.txt"},
		},
		{
			name: "No selection made",
			mockFn: func(p survey.Prompt, response interface{}, options ...survey.AskOpt) error {
				// Simulate that the user presses Enter without selecting anything.
				if sel, ok := response.(*[]string); ok {
					*sel = []string{}
				}
				return nil
			},
			inputFiles:  []string{"file1.txt", "file2.txt"},
			expectedErr: "no files selected",
		},
		{
			name: "Error during prompting",
			mockFn: func(p survey.Prompt, response interface{}, options ...survey.AskOpt) error {
				// Simulate an error (for example, I/O error)
				return errors.New("simulated prompt failure")
			},
			inputFiles:  []string{"file1.txt", "file2.txt"},
			expectedErr: "simulated prompt failure",
		},
	}

	// Iterate through test cases.
	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Override askOne with the version specified by the test case.
			askOne = tc.mockFn

			selected, err := SimpleSelector(tc.inputFiles)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr,
					fmt.Sprintf("expected error %q to contain %q", err.Error(), tc.expectedErr))
				assert.Nil(t, selected)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, selected)
			}
		})
	}
}

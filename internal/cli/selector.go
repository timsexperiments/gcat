package cli

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

// askOne is a variable that points to the survey.AskOne function. In production
// it will behave as usual, but tests can override it.
var askOne = survey.AskOne

// SimpleSelector uses Survey’s MultiSelect to prompt the user.
func SimpleSelector(files []string) ([]string, error) {
	sort.Strings(files)

	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select files:",
		Options: files,
	}

	if err := askOne(prompt, &selected); err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no files selected")
	}

	return selected, nil
}

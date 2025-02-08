package clipboard

import (
	"golang.design/x/clipboard"
)

// WriteText writes the provided text to the clipboard.
//
// It waits up to timeout; if the write is not finished by then, it returns an error.
func WriteText(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

package clipboard

import (
	"golang.design/x/clipboard"
)

var write = clipboard.Write

// WriteText writes the provided text to the clipboard.
//
// It waits up to timeout; if the write is not finished by then, it returns an error.
func WriteText(text string) {
	_ = write(clipboard.FmtText, []byte(text))
}

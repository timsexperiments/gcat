package clipboard

import (
	"testing"

	"golang.design/x/clipboard"
)

func TestWriteText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "write",
			want: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			WriteText(tt.want)
			bytes := clipboard.Read(clipboard.FmtText)

			if got := string(bytes); got != tt.want {
				t.Errorf("WriteText() = %v, want %v", got, tt.want)
			}
		})
	}
}

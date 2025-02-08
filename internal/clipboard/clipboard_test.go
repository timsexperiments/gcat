package clipboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			var answer string
			write = func(_ clipboard.Format, text []byte) <-chan struct{} {
				answer = string(text)
				return nil
			}

			WriteText(tt.want)

			assert.Equal(t, tt.want, answer)
		})
	}
}

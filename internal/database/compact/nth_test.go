package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNth(t *testing.T) {
	t.Parallel()

	const three = "a\nb\nc"

	tests := []struct {
		name   string
		s      string
		i      int
		want   string
		wantOK bool
	}{
		{"first", three, 0, "a", true},
		{"middle", three, 1, "b", true},
		{"last", three, 2, "c", true},
		{"past end", three, 3, "", false},
		{"negative", three, -1, "", false},
		{"single field", "solo", 0, "solo", true},
		{"single field past end", "lone", 1, "", false},
		{"empty string", "", 0, "", true},
		{"empty string past end", "", 1, "", false},
		{"empty fields", "\n\n", 1, "", true},
		{"trailing newline", "a\n", 1, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := nth([]byte(tt.s), tt.i)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

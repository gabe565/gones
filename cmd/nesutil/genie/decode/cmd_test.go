package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_decode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    decodeResult
		wantErr require.ErrorAssertionFunc
	}{
		{"upper", args{"YEUZUGAA"}, decodeResult{
			Code:    "YEUZUGAA",
			Address: 0xACB3,
			Replace: 0x07,
			Compare: 0x00,
		}, require.NoError},
		{"lower", args{"yeuzugaa"}, decodeResult{
			Code:    "YEUZUGAA",
			Address: 0xACB3,
			Replace: 0x07,
			Compare: 0x00,
		}, require.NoError},
		{"valid YELZUGAA", args{"YELZUGAA"}, decodeResult{
			Code:    "YELZUGAA",
			Address: 0xACB3,
			Replace: 0x07,
			Compare: 0x00,
		}, require.NoError},
		{"valid SXIOPO", args{"SXIOPO"}, decodeResult{
			Code:    "SXIOPO",
			Address: 0x91D9,
			Replace: 0xAD,
			Compare: -1,
		}, require.NoError},
		{"valid SXSOPO", args{"SXSOPO"}, decodeResult{
			Code:    "SXSOPO",
			Address: 0x91D9,
			Replace: 0xAD,
			Compare: -1,
		}, require.NoError},
		{"invalid len", args{"YEUZUGA"}, decodeResult{Code: "YEUZUGA", Compare: -1}, require.Error},
		{"invalid chars", args{"YEUZUGAF"}, decodeResult{Code: "YEUZUGAF", Compare: -1}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode(tt.args.code)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

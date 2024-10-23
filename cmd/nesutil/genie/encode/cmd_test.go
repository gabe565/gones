package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_encode(t *testing.T) {
	type args struct {
		address int
		replace int
		compare int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"valid YEUZUGAA", args{
			address: 0xACB3,
			replace: 0x07,
			compare: 0x00,
		}, "YEUZUGAA", require.NoError},
		{"valid YEUZUGAA", args{
			address: 0x2CB3,
			replace: 0x07,
			compare: 0x00,
		}, "YEUZUGAA", require.NoError},
		{"valid SXIOPO", args{
			address: 0x91D9,
			replace: 0xAD,
			compare: -1,
		}, "SXIOPO", require.NoError},
		{"valid SXIOPO", args{
			address: 0x11D9,
			replace: 0xAD,
			compare: -1,
		}, "SXIOPO", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encode(tt.args.address, tt.args.replace, tt.args.compare)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

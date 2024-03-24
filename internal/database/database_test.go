package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNameByHash(t *testing.T) {
	t.Parallel()
	type args struct {
		md5 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"Mario 3", args{"85f0ddddfe4ab67c42aba48498f42fdc"}, "Super Mario Bros. 3 (USA)", assert.NoError},
		{"Metroid", args{"397d10e475266ad28144a5fa6ec3c466"}, "Metroid (USA)", assert.NoError},
		{"Not Found", args{"a"}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindNameByHash(tt.args.md5)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

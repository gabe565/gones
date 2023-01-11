package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindNameByHash(t *testing.T) {
	type args struct {
		md5 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Mario 3", args{"85f0ddddfe4ab67c42aba48498f42fdc"}, "Super Mario Bros. 3 (USA)", false},
		{"Metroid", args{"397d10e475266ad28144a5fa6ec3c466"}, "Metroid (USA)", false},
		{"Not Found", args{"a"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindNameByHash(tt.args.md5)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

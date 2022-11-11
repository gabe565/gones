package cartridge

import "testing"

func Test_hasBattery(t *testing.T) {
	type args struct {
		data byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"false", args{}, false},
		{"true", args{0b10}, true},
		{"extraneous true", args{0b1011}, true},
		{"extraneous false", args{0b1001}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasBattery(tt.args.data); got != tt.want {
				t.Errorf("hasBattery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMirror(t *testing.T) {
	type args struct {
		data byte
	}
	tests := []struct {
		name string
		args args
		want Mirror
	}{
		{"horizontal", args{}, Horizontal},
		{"vertical", args{0b1}, Vertical},
		{"four screen", args{0b1001}, FourScreen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMirror(tt.args.data); got != tt.want {
				t.Errorf("getMirror() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMapper(t *testing.T) {
	type args struct {
		data [2]byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{"0", args{}, 0},
		{"1", args{[2]byte{0b10000}}, 1},
		{"2", args{[2]byte{0b100000}}, 2},
		{"40", args{[2]byte{0b10000000, 0b100000}}, 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMapper(tt.args.data); got != tt.want {
				t.Errorf("getMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

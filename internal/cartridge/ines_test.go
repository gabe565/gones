package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_INESFileHeader_Battery(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [10]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{}, false},
		{"true", fields{[10]byte{2}}, true},
		{"extraneous true", fields{[10]byte{0xff}}, true},
		{"extraneous false", fields{[10]byte{0xff ^ 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			if got := i.Battery(); got != tt.want {
				t.Errorf("Battery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_INESFileHeader_Mirror(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [10]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   Mirror
	}{
		{"horizontal", fields{}, Horizontal},
		{"vertical", fields{[10]byte{1}}, Vertical},
		{"four screen", fields{[10]byte{0x8 | 0x1}}, FourScreen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			if got := i.Mirror(); got != tt.want {
				t.Errorf("Mirror() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_INESFileHeader_Mapper(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [10]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"0", fields{}, 0},
		{"1", fields{[10]byte{0x10}}, 1},
		{"2", fields{[10]byte{0x20}}, 2},
		{"40", fields{[10]byte{0x80, 0x20}}, 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			if got := i.Mapper(); got != tt.want {
				t.Errorf("Mapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestINESFileHeader_SetBattery(t *testing.T) {
	t.Parallel()
	type fields struct {
		Control [10]byte
	}
	type args struct {
		v bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   byte
	}{
		{"false to true", fields{}, args{true}, 2},
		{"true to false", fields{Control: [10]byte{2}}, args{false}, 0},
		{"false to false", fields{}, args{false}, 0},
		{"true to true", fields{Control: [10]byte{2}}, args{true}, 2},
		{"extraneous unchanged when true", fields{Control: [10]byte{0xff}}, args{true}, 0xff},
		{"extraneous unchanged when false", fields{Control: [10]byte{0xff}}, args{false}, 0xfd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			i.SetBattery(tt.args.v)
			assert.Equal(t, tt.args.v, i.Battery())
			assert.Equal(t, tt.want, i.Control[0])
		})
	}
}

func TestINESFileHeader_SetMirror(t *testing.T) {
	t.Parallel()
	type fields struct {
		Control [10]byte
	}
	type args struct {
		v Mirror
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   byte
	}{
		{"horizontal to horizontal", fields{}, args{Horizontal}, 0},
		{"horizontal to vertical", fields{}, args{Vertical}, 1},
		{"vertical to horizontal", fields{Control: [10]byte{1}}, args{Vertical}, 1},
		{"vertical to vertical", fields{Control: [10]byte{1}}, args{Vertical}, 1},
		{"horizontal to four screen", fields{Control: [10]byte{1}}, args{FourScreen}, 0x8},
		{"extraneous unchanged when horizontal", fields{Control: [10]byte{0xff}}, args{Horizontal}, 0xf6},
		{"extraneous unchanged when vertical", fields{Control: [10]byte{0xff}}, args{Vertical}, 0xf7},
		{"extraneous unchanged when four screen", fields{Control: [10]byte{0xff}}, args{FourScreen}, 0xfe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			i.SetMirror(tt.args.v)
			assert.Equal(t, tt.args.v, i.Mirror())
			assert.Equal(t, tt.want, i.Control[0])
		})
	}
}

func TestINESFileHeader_SetMapper(t *testing.T) {
	t.Parallel()
	type fields struct {
		Control [10]byte
	}
	type args struct {
		v uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [10]byte
	}{
		{"0 to 1", fields{}, args{1}, [10]byte{0x10, 0}},
		{"0 to 71", fields{}, args{71}, [10]byte{0x70, 0x40}},
		{"71 to 0", fields{[10]byte{0x70, 0x40}}, args{0}, [10]byte{}},
		{"extraneous unchanged", fields{Control: [10]byte{0xff, 0xff, 0xff}}, args{0}, [10]byte{0xf, 0xf, 0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := INESFileHeader{Control: tt.fields.Control}
			i.SetMapper(tt.args.v)
			assert.Equal(t, tt.args.v, i.Mapper())
			assert.Equal(t, tt.want, i.Control)
		})
	}
}

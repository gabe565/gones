package cartridge

import "testing"

func Test_iNESFileHeader_Battery(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [2]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"false", fields{}, false},
		{"true", fields{[2]byte{0b10, 0}}, true},
		{"extraneous true", fields{[2]byte{0b1011, 0}}, true},
		{"extraneous false", fields{[2]byte{0b1001, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := iNESFileHeader{
				Control: tt.fields.Control,
			}
			if got := i.Battery(); got != tt.want {
				t.Errorf("Battery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iNESFileHeader_Mirror(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [2]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   Mirror
	}{
		{"horizontal", fields{}, Horizontal},
		{"vertical", fields{[2]byte{0b1, 0}}, Vertical},
		{"four screen", fields{[2]byte{0b1001, 0}}, FourScreen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := iNESFileHeader{
				Control: tt.fields.Control,
			}
			if got := i.Mirror(); got != tt.want {
				t.Errorf("Mirror() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iNESFileHeader_Mapper(t *testing.T) {
	t.Parallel()

	type fields struct {
		Control [2]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"0", fields{}, 0},
		{"1", fields{[2]byte{0b10000}}, 1},
		{"2", fields{[2]byte{0b100000}}, 2},
		{"40", fields{[2]byte{0b10000000, 0b100000}}, 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := iNESFileHeader{
				Control: tt.fields.Control,
			}
			if got := i.Mapper(); got != tt.want {
				t.Errorf("Mapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

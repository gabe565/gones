package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_blarggAPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"len ctr", "roms/apu_test/rom_singles/1-len_ctr.nes", 0, "1-len_ctr\n\nPassed"},
		{"len table", "roms/apu_test/rom_singles/2-len_table.nes", 0, "2-len_table\n\nPassed"},
		{"IRQ flag", "roms/apu_test/rom_singles/3-irq_flag.nes", 0, "3-irq_flag\n\nPassed"},
		{"jitter", "roms/apu_test/rom_singles/4-jitter.nes", 2, "Frame irq is set too soon\n\n4-jitter\n\nFailed #2"},
		{
			"len timing",
			"roms/apu_test/rom_singles/5-len_timing.nes",
			2,
			"Channel: 0\n\nFirst length of mode 0 is too soon\n\n5-len_timing\n\nFailed #2",
		},
		{
			"IRQ flag timing",
			"roms/apu_test/rom_singles/6-irq_flag_timing.nes",
			2,
			"Flag first set too soon\n\n6-irq_flag_timing\n\nFailed #2",
		},
		{"DMC basics", "roms/apu_test/rom_singles/7-dmc_basics.nes", 0, "7-dmc_basics\n\nPassed"},
		{"DMC rates", "roms/apu_test/rom_singles/8-dmc_rates.nes", 0, "8-dmc_rates\n\nPassed"},
		{"reset clears $4015", "roms/apu_reset/4015_cleared.nes", 0, "4015_cleared\n\nPassed"},
		{"reset clears IRQ", "roms/apu_reset/irq_flag_cleared.nes", 0, "irq_flag_cleared\n\nPassed"},
		{
			"reset $4017 timing",
			"roms/apu_reset/4017_timing.nes",
			3,
			"Delay after effective $4017 write: 0\n\nFrame IRQ flag should be set sooner after power/reset\n\n4017_timing\n\nFailed #3",
		},
		{
			"reset $4017 written",
			"roms/apu_reset/4017_written.nes",
			2,
			"At power, $4017 should be written with $00\n\n4017_written\n\nFailed #2",
		},
		{
			"reset len ctrs enabled",
			"roms/apu_reset/len_ctrs_enabled.nes",
			2,
			"At power, length counters should be enabled\n\nlen_ctrs_enabled\n\nFailed #2",
		},
		{"reset works immediately", "roms/apu_reset/works_immediately.nes", 0, "works_immediately\n\nPassed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, msgTypeSRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.Equal(t, tt.wantStatus, getBlarggStatus(test))
			assert.Equal(t, tt.want, getBlarggMessage(test, msgTypeSRAM))
		})
	}
}

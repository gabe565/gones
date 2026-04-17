package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blarggCPUTimingWant = `6502 TIMING TEST (16 SECONDS)
OFFICIAL INSTRUCTIONS ONLY
PASSED`

func Test_blarggCPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		msgType    msgType
		wantStatus status
		want       string
	}{
		{"instructions", "roms/instr_test-v5/all_instrs.nes", msgTypeSRAM, 0, "All 16 tests passed"},
		{"instruction timing", "roms/cpu_timing_test6/cpu_timing_test.nes", msgTypePPUVRAM, -1, blarggCPUTimingWant},
		{
			"branch timing basics",
			"roms/branch_timing_tests/1.Branch_Basics.nes",
			msgTypePPUVRAM,
			-1,
			"BRANCH TIMING BASICS\nPASSED",
		},
		{
			"branch timing backward",
			"roms/branch_timing_tests/2.Backward_Branch.nes",
			msgTypePPUVRAM,
			-1,
			"BACKWARD BRANCH TIMING\nPASSED",
		},
		{
			"branch timing forward",
			"roms/branch_timing_tests/3.Forward_Branch.nes",
			msgTypePPUVRAM,
			-1,
			"FORWARD BRANCH TIMING\nPASSED",
		},
		{"ram after reset", "roms/cpu_reset/ram_after_reset.nes", msgTypeSRAM, 0, "ram_after_reset\n\nPassed"},
		{
			"registers after reset",
			"roms/cpu_reset/registers.nes",
			msgTypeSRAM,
			0,
			"A  X  Y  P  S\n34 56 78 FF 0F \n\nregisters\n\nPassed",
		},
		{
			"instr misc abs x wrap",
			"roms/instr_misc/rom_singles/01-abs_x_wrap.nes",
			msgTypeSRAM,
			0,
			"01-abs_x_wrap\n\nPassed",
		},
		{
			"instr misc branch wrap",
			"roms/instr_misc/rom_singles/02-branch_wrap.nes",
			msgTypeSRAM,
			0,
			"02-branch_wrap\n\nPassed",
		},
		{
			"instr misc dummy reads",
			"roms/instr_misc/rom_singles/03-dummy_reads.nes",
			msgTypeSRAM,
			0,
			"03-dummy_reads\n\nPassed",
		},
		{
			"instr misc dummy reads APU",
			"roms/instr_misc/rom_singles/04-dummy_reads_apu.nes",
			msgTypeSRAM,
			0,
			"04-dummy_reads_apu\n\nPassed",
		},
		{
			"interrupts cli latency",
			"roms/cpu_interrupts_v2/rom_singles/1-cli_latency.nes",
			msgTypeSRAM,
			5,
			"CLI SEI should allow only one IRQ just after SEI\n\n1-cli_latency\n\nFailed #5",
		},
		{
			"interrupts NMI and BRK",
			"roms/cpu_interrupts_v2/rom_singles/2-nmi_and_brk.nes",
			msgTypeSRAM,
			1,
			"NMI BRK 00\n27  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n26  36  00 \n\nCB0122C0\n2-nmi_and_brk\n\nFailed",
		},
		{
			"interrupts NMI and IRQ",
			"roms/cpu_interrupts_v2/rom_singles/3-nmi_and_irq.nes",
			msgTypeSRAM,
			1,
			"NMI BRK\n23  00 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n27  23 \n\nE589905A\n3-nmi_and_irq\n\nFailed",
		},
		{
			"interrupts IRQ and DMA",
			"roms/cpu_interrupts_v2/rom_singles/4-irq_and_dma.nes",
			msgTypeSRAM,
			1,
			"0 +0\n1 +1\n1 +2\n1 +3\n1 +4\n1 +5\n2 +6\n2 +7\n4 +8\n4 +9\n7 +10\n7 +11\n8 +12\n8 +13\n...\n8 +524\n8 +525\n8 +526\n8 +527\n\n2343931E\n4-irq_and_dma\n\nFailed",
		},
		{
			"interrupts branch delays IRQ",
			"roms/cpu_interrupts_v2/rom_singles/5-branch_delays_irq.nes",
			msgTypeSRAM,
			1,
			"test_jmp\nT+ CK PC\n00 05 04 \n01 04 04 \n02 03 04 \n03 02 04 \n04 01 04 \n05 03 07 \n06 02 07 \n07 03 08 \n08 02 08 \n09 01 08 \n\n\n74DC0BFA\n5-branch_delays_irq\n\nFailed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, tt.msgType)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.Equal(t, tt.wantStatus, getBlarggStatus(test))
			assert.Equal(t, tt.want, getBlarggMessage(test, tt.msgType))
		})
	}
}

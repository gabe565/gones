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
			"NMI BRK\n23  00 \n21  00 \n21  00 \n20  00 \n20  00 \n20  00 \n20  00 \n00  00 \n00  00 \n00  00 \n00  00 \n00  00 \n\n7A096051\n3-nmi_and_irq\n\nFailed",
		},
		{
			"interrupts IRQ and DMA",
			"roms/cpu_interrupts_v2/rom_singles/4-irq_and_dma.nes",
			msgTypeSRAM,
			1,
			"53 +0\n53 +1\n53 +2\n53 +3\n53 +4\n53 +5\n53 +6\n53 +7\n53 +8\n53 +9\n53 +10\n53 +11\n53 +12\n53 +13\n...\n53 +524\n53 +525\n53 +526\n53 +527\n\nD927EAD0\n4-irq_and_dma\n\nFailed",
		},
		{
			"interrupts branch delays IRQ",
			"roms/cpu_interrupts_v2/rom_singles/5-branch_delays_irq.nes",
			msgTypeSRAM,
			1,
			"test_jmp\nT+ CK PC\n00 02 02 \n01 01 02 \n02 02 00 \n03 02 00 \n04 02 02 \n05 01 02 \n06 01 00 \n07 01 04 \n08 02 04 \n09 02 00 \n\n\n663F1536\n5-branch_delays_irq\n\nFailed",
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

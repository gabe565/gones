# GoNES

[![Build](https://github.com/gabe565/gones/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/gones/actions/workflows/build.yml)

GoNES is another NES emulator written in Go.

This emulator is in the early development stages and is not yet stable.

I know that there are plenty of great NES emulators, but I love writing Go,
and I wanted to learn more about the NES and emulation in general.

## Install

### Binary

Eventually, binaries will be attached to releases.
For now, binaries can be downloaded from CI build artifacts.
1. Go to [the build action](https://github.com/gabe565/gones/actions/workflows/build.yml).
2. Click the latest build job.
3. Scroll down to "Artifacts" and download `dist`.
4. Run the binary for your OS in the downloaded zip file!

### From Source

Make sure you have [Go](https://go.dev/doc/install) and the [requirements](#requirements) installed, then run:

```shell
go install github.com/gabe565/gones
```

## Requirements

Rendering uses [faiface/pixel](https://github.com/faiface/pixel) which requires
OpenGL development libraries to compile.
See [pixel requirements](https://github.com/faiface/pixel#requirements).

## Usage

```shell
gones ROM_FILE
```

## Keybinds

Keys are currently hardcoded in [`internal/controller/keymap.go`](./internal/controller/keymap.go).
Eventually, this will be configurable in the UI.

### Player 1

| Button     | Key         |
|------------|-------------|
| A          | A           |
| B          | S           |
| Directions | Arrow Keys  |
| Start      | Enter       |
| Select     | Right Shift |

### Player 2

Player 2 controller is implemented, but still undecided on keybinds.

## Milestones

- [x] CPU implementation
  - CPU is stable, and `nestest.nes` passes.
- [x] Cartridge implementation
  - [ ] Common mappers implemented
    - Currently only mapper 0 is supported.
- [x] PPU implementation (graphics)
  - [x] Background rendering 
  - [x] Sprite rendering
  - PPU emulation is still rough. Games that rely on certain quirks of the
    original PPU will not work.
- [x] GUI
  - Video output implemented using [faiface/pixel](https://github.com/faiface/pixel), but menu options need to be added.
- [x] Basic controller support
  - [x] Player 1
  - [x] Player 2 (implemented, but need to decide on keybinds)
- [ ] APU implementation (audio)
  - Planning to use [faiface/beep](https://github.com/faiface/beep)
- [ ] Save games
- [ ] Save states
- [ ] Preferences (remap controllers, video config, sound config, etc)
- [ ] Battery support
- [ ] Cheats

## References

- [NESDev wiki](https://www.nesdev.org/wiki/Nesdev_Wiki)
- [NESDev 6502 Reference](https://www.nesdev.org/obelisk-6502-guide/)
- [NESDev Undocumented Opcodes](https://www.nesdev.org/undocumented_opcodes.txt)
- [Writing an NES Emulator in Rust](https://bugzmanov.github.io/nes_ebook/)

<div align="center">
<h1>
  <img src="assets/icon.svg" width="192" align="left"/>
  GoNES
</h1>

[![Build](https://github.com/gabe565/gones/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/gones/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/gones)](https://goreportcard.com/report/github.com/gabe565/gones)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gabe565_gones&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabe565_gones)

_An NES emulator written in Go._

</div>

## Install

### Release

Automated builds are uploaded during the release process. See the [latest release](https://github.com/gabe565/gones/releases/latest) for download links.

### From Source

Make sure you have [Go](https://go.dev/doc/install) and the [requirements](#requirements) installed, then run:

```shell
go install github.com/gabe565/gones@latest
```

## Requirements

Rendering uses [hajimehoshi/ebiten](https://github.com/hajimehoshi/ebiten) which requires
OpenGL development libraries to compile.
See [ebiten install steps](https://ebitengine.org/en/documents/install.html).

## Usage

In a terminal, run:

```shell
gones ROM_FILE
```

See [docs](./docs/gones.md) for the full command line usage documentation.

## Configuration

A configuration file will be generated the first time GoNES is run. Depending on your operating system, the file will be available at:
- **Windows:** `%AppData%\gones\config.toml`
- **macOS:** `~/Library/Application Support/pwgen-go/config.toml`
- **Linux:** `~/.config/gones/config.toml`

An example configuration is also available at [`config_example.toml`](config_example.toml).

## Keybinds

Keys are configurable, but the default values are listed below.

### Player 1

| Nintendo   | Emulator    |
|------------|-------------|
| A          | M           |
| B          | N           |
| Directions | WASD        |
| Start      | Enter       |
| Select     | Right Shift |
| A (Turbo)  | K           |
| B (Turbo)  | J           |

### Player 2

<details>
  <summary>Click to expand</summary>

| Nintendo   | Emulator          |
|------------|-------------------|
| A          | Num Pad 3         |
| B          | Num Pad 2         |
| Directions | Home/Del/End/PgDn |
| Start      | Num Pad Enter     |
| Select     | Num Pad Plus      |
| A (Turbo)  | Num Pad 6         |
| B (Turbo)  | Num Pad 5         |

</details>

### Other

| Action            | Key      |
|-------------------|----------|
| Save State        | F1       |
| Load State        | F5       |
| Undo Save State   | Shift+F1 |
| Undo Load State   | Shift+F5 |
| Fast Forward      | F (Hold) |
| Reset             | R (Hold) |
| Toggle Fullscreen | F11      |

#### Debugging

<details>
  <summary>Click to expand</summary>

| Action                                            | Key |
|---------------------------------------------------|-----|
| Toggle step debugging                             | `   |
| Toggle stdout trace log (when step debug enabled) | Tab |
| Step to next frame                                | 1   |
| Run to next render                                | 2   |

</details>

## Milestones

- [x] CPU implementation
  - CPU is stable, and `nestest.nes` passes.
- [x] Cartridge implementation
  - [x] Support for mappers
  - [x] Common mappers implemented
    - Supported mappers: 0, 1, 2, 3, 4, 7, 69, 71 (84.34% of official NES games)
- [x] PPU implementation (graphics)
  - [x] Background rendering
  - [x] Sprite rendering
- [x] GUI
  - Rendering works, but menu options need to be added.
- [x] Basic controller support
  - [x] Player 1
  - [x] Player 2
  - [ ] External controllers
- [x] APU implementation (audio)
- [x] Save file for games with batteries
- [x] Save states
- [x] Configuration (remap controllers, video config, sound config, etc)
  - [x] Config file
  - [ ] Config UI
- [ ] Cheats

## References

- [NESDev wiki](https://www.nesdev.org/wiki/Nesdev_Wiki)
- [NESDev 6502 Reference](https://www.nesdev.org/obelisk-6502-guide/)
- [NESDev Undocumented Opcodes](https://www.nesdev.org/undocumented_opcodes.txt)
- [Writing an NES Emulator in Rust](https://bugzmanov.github.io/nes_ebook/)
- [fogleman/nes](https://github.com/fogleman/nes)
- [ltriant/nes](https://github.com/ltriant/nes)
- [i82orbom/nesgo](https://github.com/i82orbom/nesgo)
- [No-Intro](https://no-intro.org)
- [NES Composite Palette Project](https://www.firebrandx.com/nespalette.html)

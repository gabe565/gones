package callbacks

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/cpu"
	"image"
	"image/color"
	"math/rand"
	"time"
)

func Snake(win *pixelgl.Window) Callback {
	rand.Seed(time.Now().UnixNano())

	frameStart := 0x200
	frameEnd := 0x600
	frame := make([]byte, frameEnd-frameStart+1)

	return func(c *cpu.CPU) error {
		var drawFrame bool
		for k, prevVal := range frame {
			val := c.MemRead(uint16(frameStart + k))
			if !drawFrame && val != prevVal {
				drawFrame = true
			}
			frame[k] = val
		}

		if win.JustPressed(pixelgl.KeyEscape) {
			return cpu.ErrBrk
		} else if win.JustPressed(pixelgl.KeyW) {
			c.MemWrite(0xFF, 0x77)
		} else if win.JustPressed(pixelgl.KeyA) {
			c.MemWrite(0xFF, 0x61)
		} else if win.JustPressed(pixelgl.KeyS) {
			c.MemWrite(0xFF, 0x73)
		} else if win.JustPressed(pixelgl.KeyD) {
			c.MemWrite(0xFF, 0x64)
		}

		if drawFrame {
			img := image.NewRGBA(image.Rect(0, 0, 32, 32))
			for k, pxl := range frame {
				var c color.Color
				switch pxl {
				case 0x0: // Black
					c = color.Black
				case 0x1: // White
					c = color.White
				case 0x2: // Red
					c = color.RGBA{R: 0xFF, A: 0xFF}
				case 0x3: // Cyan
					c = color.RGBA{G: 0xFF, B: 0xFF, A: 0xFF}
				case 0x4: // Purple
					c = color.RGBA{R: 0x80, B: 0x80, A: 0xFF}
				case 0x5: // Green
					c = color.RGBA{G: 0xFF, A: 0xFF}
				case 0x6: // Blue
					c = color.RGBA{B: 0xFF, A: 0xFF}
				case 0x7: // Yellow
					c = color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF}
				case 0x8: // Orange
					c = color.RGBA{R: 0xFF, G: 0xA5, A: 0xFF}
				case 0x9: // Brown
					c = color.RGBA{R: 0xA5, G: 0x2A, B: 0x2A, A: 0xFF}
				case 0xA: // Light red
					c = color.RGBA{R: 0xFF, G: 0x4F, B: 0x4D, A: 0xFF}
				case 0xB: // Dark gray
					c = color.Gray{Y: 0xA9}
				case 0xC: // Grey
					c = color.Gray{Y: 0x80}
				case 0xD: // Light green
					c = color.RGBA{R: 0x90, G: 0xEE, B: 0x90, A: 0xFF}
				case 0xE: // Light blue
					c = color.RGBA{R: 0xAD, G: 0xD8, B: 0xE6, A: 0xFF}
				case 0xF: // Light gray
					c = color.Gray{Y: 0xD3}
				default:
					c = color.White
				}
				img.Set(k%32, k/32, c)
			}

			win.Clear(color.Black)
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 10))
			win.Update()
			time.Sleep(30 * time.Millisecond)
		}

		seed := byte(rand.Intn(15) + 1)
		c.MemWrite(0xFE, seed)

		return nil
	}
}

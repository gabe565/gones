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

		if drawFrame {
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

			img := image.NewRGBA(image.Rect(0, 0, 32, 32))
			for k, pxl := range frame {
				var c color.Color
				switch pxl {
				case 0:
					c = color.Black
				case 1:
					c = color.White
				case 2, 9:
					c = color.Gray16{Y: 0x8888}
				case 3, 10:
					c = color.RGBA{R: 0xFF, A: 0xFF}
				case 4, 11:
					c = color.RGBA{G: 0xFF, A: 0xFF}
				case 5, 12:
					c = color.RGBA{B: 0xFF, A: 0xFF}
				case 6, 13:
					c = color.RGBA{R: 0xFF, B: 0xFF, A: 0xFF}
				case 7, 14:
					c = color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF}
				default:
					c = color.RGBA{G: 0xFF, B: 0xFF, A: 0xFF}
				}
				img.Set(k%32, k/32, c)
			}

			win.Clear(color.Black)
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 10))
			win.Update()

			c.MemWrite(0xFE, byte(rand.Intn(15)+1))
		}

		return nil
	}
}

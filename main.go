package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/games"
	"image"
	"image/color"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, 10*32, 10*32),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := cpu.New()
	c.PrgRomAddr = 0x600
	c.Load(games.Snake)
	c.Memory[0xFF] = 0x77
	c.Reset()

	var lastImg image.Image
	c.Callback = func(c *cpu.CPU) {
		img := image.NewRGBA(image.Rect(0, 0, 32, 32))
		for k, pxl := range c.Memory[0x200:0x600] {
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

		if !reflect.DeepEqual(img, lastImg) {
			if win.JustPressed(pixelgl.KeyEscape) {
				os.Exit(0)
			} else if win.JustPressed(pixelgl.KeyW) {
				c.Memory[0xFF] = 0x77
			} else if win.JustPressed(pixelgl.KeyA) {
				c.Memory[0xFF] = 0x61
			} else if win.JustPressed(pixelgl.KeyS) {
				c.Memory[0xFF] = 0x73
			} else if win.JustPressed(pixelgl.KeyD) {
				c.Memory[0xFF] = 0x64
			}

			win.Clear(color.Black)
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 10))
			win.Update()
			lastImg = img
		}
		time.Sleep(70000 * time.Nanosecond)

		c.Memory[0xFE] = uint8(rand.Int31())
	}

	if err := c.Run(); err != nil {
		panic(err)
	}
}

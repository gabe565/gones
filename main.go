package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/cpu"
	"image"
	"image/color"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func main() {
	if err := NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func run(path string) error {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, 10*32, 10*32),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	c, err := console.New(path)
	if err != nil {
		return err
	}
	c.Reset()

	rand.Seed(time.Now().UnixNano())

	var lastImg image.Image
	c.Callback = func(c *cpu.CPU) {
		if win.Pressed(pixelgl.KeyEscape) {
			os.Exit(0)
		} else if win.Pressed(pixelgl.KeyW) {
			c.MemWrite(0xFF, 0x77)
		} else if win.Pressed(pixelgl.KeyA) {
			c.MemWrite(0xFF, 0x61)
		} else if win.Pressed(pixelgl.KeyS) {
			c.MemWrite(0xFF, 0x73)
		} else if win.Pressed(pixelgl.KeyD) {
			c.MemWrite(0xFF, 0x64)
		}

		img := image.NewRGBA(image.Rect(0, 0, 32, 32))
		for addr := 0x200; addr <= 0x600; addr += 1 {
			pxl := c.MemRead(uint16(addr))
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
			k := addr - 0x200
			img.Set(k%32, k/32, c)
		}

		if !reflect.DeepEqual(img, lastImg) {
			win.Clear(color.Black)
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 10))
			win.Update()
			lastImg = img
			time.Sleep(time.Second / 60)
		}

		c.MemWrite(0xFE, byte(rand.Intn(15)+1))
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

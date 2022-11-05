package main

import (
	"fmt"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/games"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"reflect"
)

func main() {
	c := cpu.New()
	c.PrgRomAddr = 0x600
	c.Load(games.Snake)
	c.Reset()

	var count uint
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
			count += 1

			resized := image.NewRGBA(image.Rect(0, 0, 10*32, 10*32))
			draw.NearestNeighbor.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)

			f, err := os.Create(fmt.Sprintf("frame_%d.png", count))
			if err != nil {
				panic(err)
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			if err := png.Encode(f, resized); err != nil {
				panic(err)
			}
			lastImg = img
		}

		c.Memory[0xFE] = uint8(rand.Int31())
		c.Memory[0xFF] = 0x61
	}

	if err := c.Run(); err != nil {
		panic(err)
	}
}

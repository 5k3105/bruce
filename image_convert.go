package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

/*
disp_buf_field_len: [13]
debug_width:        [384]
debug_height:       [247]
x_inner:            [0]
y_inner:            [0]
width_inner:        [384]
height_inner:       [247]
bits:               [8]
buffer_len:         [94848]

	palette := palette_get()
	fmt.Println(palette)

*/

func raw_to_png(path, filename string, data []byte) {
	data = append(data, []byte{0, 0, 0, 0}...)
	width := 384
	height := 247
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cs := colorset[data[i]]
			img.Set(x, y, color.RGBA{cs[0], cs[1], cs[2], 0xff})
			i++
		}
	}

	f, _ := os.Create(path + filename + ".png")
	png.Encode(f, img)
}

var colorset = map[byte][3]byte{
	/// Black
	0x00: [3]byte{0x00, 0x00, 0x00},

	/// White
	0x01: [3]byte{0xFF, 0xFF, 0xFF},

	/// Red
	0x02: [3]byte{0x7C, 0x35, 0x2B},

	/// Cyan
	0x03: [3]byte{0x5A, 0xA6, 0xB1},

	/// Purple
	0x04: [3]byte{0x69, 0x41, 0x85},

	/// Green
	0x05: [3]byte{0x5D, 0x86, 0x43},

	/// Blue
	0x06: [3]byte{0x21, 0x2E, 0x78},

	/// Yellow
	0x07: [3]byte{0xCF, 0xBE, 0x6F},

	/// Orange
	0x08: [3]byte{0x89, 0x4A, 0x26},

	/// Brown
	0x09: [3]byte{0x5B, 0x33, 0x0},

	/// Light Red
	0x0a: [3]byte{0xAF, 0x64, 0x59},

	/// Dark Gray
	0x0b: [3]byte{0x43, 0x43, 0x43},

	/// Medium Gray
	0x0c: [3]byte{0x6B, 0x6B, 0x6B},

	/// Light Green
	0x0d: [3]byte{0xA0, 0xCB, 0x84},

	/// Light Blue
	0x0e: [3]byte{0x56, 0x65, 0xB3},

	/// Light Gray
	0x0f: [3]byte{0x95, 0x95, 0x95},
}

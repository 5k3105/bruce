package main

import (
	"io/ioutil"
	"strconv"
)

type level_record struct {
	compressed []byte
	screen     [1000]byte
	color      [1000]byte
	lanterns   []byte
	code       []byte
}

var levelmaps map[int]level_record

func load_() { /// colors not decompressed
	load_level_data("./original_data/", "LevelData")
	lvl, _ := strconv.Atoi(level)
	load_level_mem(lvl)
}

/*
func load_level_mem(lvl int) {
	for inc, data := range levelmaps[lvl].screen {
		mem_ram_[screen_mem_+uint16(inc)] = data
	}

	for inc, data := range levelmaps[lvl].color {
		mem_ram_[color_mem_+uint16(inc)] = data
	}
}
*/

func load_level_mem(lvl int) {
	load_file_memory(map_path, "cset1", 0x3800)
	data := levelmaps[lvl].screen
	send_command(memory_set(0x0478, data[:880]))
	data = levelmaps[lvl].color
	send_command(memory_set(0xd878, data[:880]))
}

func load_level_data(path, filename string) {
	file, err := ioutil.ReadFile(path + filename)
	if err != nil {
		println(err.Error)
	}
	parse_levels(file)
}

func parse_levels(bytes []byte) {
	var (
		r                  bool
		ccol, cols         uint16 = 40, 40
		clevel             int
		soffs, boffs, offs uint16
		x, c, b            uint8
	)

	var screen, color = [1000]byte{}, [1000]byte{}
	levelmaps = make(map[int]level_record)

top:
	if int(boffs) > len(bytes) {
		return
	}

	b = bytes[boffs]      /// 16 8D
	if !ISSET_BIT(b, 7) { /// bpl
		r = true
		goto init
	}
	if b == 0x80 {
		goto next_level
	}
	b &= 0x7F

	r = false

init:
	x = b
	boffs++

data_read:
	b = bytes[boffs] /// 8A
	if !ISSET_BIT(b, 7) {
		goto write_screen
	}
	c = 0xF8
	color[soffs] = c
	color[soffs+cols] = c
	b &= 0x7F

write_screen:
	screen[soffs] = b
	screen[soffs+cols] = b | 0x80

	if !r {
		boffs++
	}
	soffs++

	ccol--
	if ccol == 0 {
		ccol = cols
		soffs += cols
	}

	x--
	if x != 0 {
		goto data_read
	}
	if r {
		boffs++
	}

	goto top

next_level:
	lvl := level_record{
		compressed: bytes[offs:boffs],
		screen:     screen,
		color:      color,
	}

	levelmaps[clevel] = lvl
	soffs = 0
	offs = boffs
	screen, color = [1000]byte{}, [1000]byte{}

	clevel++
	if clevel == 20 {
		return
	}
	boffs++
	goto top
}

func ISSET_BIT(v, b uint8) bool {
	return (v & (1 << uint8(b))) != 0
}

/* level offsets
68F9
6A09
6B3E
6C9D
6E1B
6F82
70DD
7227
7311
7491
75D8
774C
7843
79C6
7A46
7BB4
7C5C
7D8D
7E90
7FE0
*/

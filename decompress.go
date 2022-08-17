package main

type level_record struct {
	compressed []byte
	screen     [1000]byte
	color      [1000]byte
	lanterns   []byte
	code       []byte
}

var levelmaps map[int]level_record

func decompress_editor_level(level string) {
	bytes := load_file(map_path, "level"+level, false)
	lrec := decompress_level(bytes)
	/// load_file_memory(map_path, "cset1", editor_cset_addr)
	data := lrec.screen
	memory_set(editor_screen_addr, data[:editor_screen_length])
	data = lrec.color
	memory_set(editor_color_addr, data[:editor_color_length])
}

func compress_editor_level(level string) []byte {
	screen := memory_get(editor_screen_addr, editor_screen_length)
	color := memory_get(editor_color_addr, editor_screen_length)
	bytes := compress_level(screen, color)
	return bytes
}

func read_lanterns_level(level string) []byte {
	screen := memory_get(editor_screen_addr, editor_screen_length)
	bytes := read_lanterns(screen)
	return bytes
}

func compress_level(screen_data, color_data []byte) []byte {
	var (
		bytes              []byte
		buffer             []byte
		sd1, sd2, cd1, cd2 uint8
		repeater           bool
		byte_offset        uint16
		count              uint8
		num_columns        uint16 = 40
		column_counter     uint16
		complete           = false
	)

	for {
		sd1 = screen_data[byte_offset]
		cd1 = color_data[byte_offset]

		byte_offset++
		column_counter++

		if column_counter == num_columns {
			column_counter = 0
			byte_offset += num_columns
		}

		if int(byte_offset) >= len(screen_data) {
			complete = true
		} else {
			sd2 = screen_data[byte_offset]
			cd2 = color_data[byte_offset]
		}

		if (sd1 != sd2 || cd1 != cd2) || complete {
			if (cd1 | 0xf0) == 0xf8 {
				sd1 |= 0x80
			}
			if repeater {
				count++
				bytes = append(bytes, count)
				bytes = append(bytes, sd1)
				repeater = false
			} else {
				buffer = append(buffer, sd1)
			}
		}

		if (sd1 == sd2 && cd1 == cd2) || complete {
			count++
			if !repeater {
				if len(buffer) != 0 {
					count = uint8(len(buffer))
					count |= 0x80
					bytes = append(bytes, count)
					bytes = append(bytes, buffer...)
					buffer = []byte{}
				}
				repeater = true
				count = 1
			}
		}

		if complete {
			bytes = append(bytes, 0x80)
			return bytes
		}

	}

	return []byte{}
}

func plus(v uint8) bool {
	return (v & (1 << uint8(7))) != 0
}

func decompress_level(bytes []byte) level_record {
	var (
		screen_data    = [1000]byte{}
		color_data     = [1000]byte{}
		color, b       uint8
		repeater       bool
		byte_offset    uint16
		screen_offset  uint16
		num_columns    uint16 = 40
		column_counter uint16
	)

	for {
		b = bytes[byte_offset]
		if b == 0x80 {
			return level_record{
				screen: screen_data,
				color:  color_data,
			}
		}

		if plus(b) {
			b &= 0x7f
			repeater = false
		} else {
			repeater = true
		}

		byte_offset++
		for i := int(b); i != 0; i-- {
			b = bytes[byte_offset]
			if plus(b) {
				b &= 0x7f
				color = 0xf8 /// black
			} else {
				color = 0xfe /// blue
			}

			screen_data[screen_offset] = b
			screen_data[screen_offset+num_columns] = b | 0x80

			color_data[screen_offset] = color
			color_data[screen_offset+num_columns] = color

			if !repeater {
				byte_offset++
			}

			screen_offset++
			column_counter++

			if column_counter == num_columns {
				column_counter = 0
				screen_offset += num_columns
			}
		}

		if repeater {
			byte_offset++
		}
	}
	return level_record{}
}

func read_lanterns(screen_data []byte) []byte {
	var (
		lantern_lookup  = []byte{0x13, 0x14, 0x31}
		lantern_convert = []byte{0x9e, 0x9f, 0x92}
		byte_offset     uint16
		bytes           = []byte{}
		b               uint8
		ok              bool
		num_columns     uint16 = 40
		column_counter  uint16
		row_counter     uint8 = 2
	)

	for {
		b = screen_data[byte_offset]
		b, ok = index_to(b, lantern_lookup, lantern_convert)
		if ok {
			bytes = append(bytes, 0x01)
			bytes = append(bytes, b)
			bytes = append(bytes, uint8(column_counter))
			bytes = append(bytes, row_counter)
		}

		byte_offset++
		column_counter++

		if column_counter == num_columns {
			column_counter = 0
			byte_offset += num_columns
			row_counter++
		}

		if int(byte_offset) >= len(screen_data) {
			bytes = append(bytes, 0xff)
			return bytes
		}

	}
}

func index_to(inp uint8, lookup, convert []byte) (uint8, bool) {
	for i, b := range lookup {
		if b == inp {
			return convert[i], true
		}
	}
	return 0, false
}

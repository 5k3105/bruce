package main

import (
	"os"
)

var (
	data_end_seq = []byte{0xff, 0xff, 0x00, 0x00} /// {0x00,0xff,0x00,0xff}
	emu_port     = "6502"                         /// x64 -binarymonitoraddress ip4://127.0.0.1:6503
	buffer_size  = 5000
	data_path    = `./data/`
	map_path     = `./map/`
	level        = ""
)

func main() {
	if len(os.Args) <= 2 {
		println("provide level number")
		return
	}
	if len(os.Args) != 3 {
		println("load (l) or save (s) ?")
		return
	}
	level = os.Args[1]
	op := os.Args[2]
	switch op {
	case "l":
		load()
	case "s":
		save()
	}
}

func load() {
	load_editor(level)
}

func save() {
	save_editor(level)
}

func load_editor(level string) {
	/// load_memory("editor.prg", 0xc000)
	load_file_memory(map_path, "cset1", 0x3800)
	load_file_memory(map_path, "map"+level, 0x0478)
	load_file_memory(map_path, "color"+level, 0xd878)
	/// jump(0xC000)
}

func save_editor(level string) {
	/**/
	screen_map := send_command(memory_get(editor_screen, editor_screen_length), RESPONSE_MEM_GET)
	save_file(map_path, "map"+level, screen_map)
	color_map := send_command(memory_get(editor_color, editor_color_length), RESPONSE_MEM_GET)
	save_file(map_path, "color"+level, color_map)

	/// how does editor store level length ? peek editor length bytes for level, lanterns before request (load labels map)
	compressed_level := send_command(memory_get(editor_compressed_level, editor_screen_length), RESPONSE_MEM_GET)
	/// bug where previous compile in editor results in smaller data length, thus cant detect compressed data length
	length := detect_data_length(compressed_level, data_end_seq) ///  average 272
	if length == 0 {
		panic("zero length: detect_data_length(compressed_level, data_end_seq)")
	}
	compressed_level = compressed_level[:length]
	save_file(data_path, "level"+level, compressed_level)

	lantern_set := send_command(memory_get(editor_lantern_set, 392), RESPONSE_MEM_GET) /// max for all levels. whats per level size?

	lantern_set = lantern_set[1:lantern_set[0]]
	save_file(data_path, "lantern"+level, lantern_set)

	println("length: compressed_level|lantern_set:", len(compressed_level), len(lantern_set))

	level_pointers, level_data := link_byte(data_path, "level", game_level_data)
	save_file(data_path, "pointers_level", level_pointers)
	save_file(data_path, "data_level", level_data)
	/**/
	lantern_pointers, lantern_data := link_word(data_path, "lantern", game_lantern_data)
	save_file(data_path, "pointers_lantern", lantern_pointers)
	save_file(data_path, "data_lantern", lantern_data)

	doors_next_level, doors_player_pos := read_doors_table()
	save_file(`./data/`, "doors_next_level", doors_next_level)
	save_file(`./data/`, "doors_player_pos", doors_player_pos)

	enemy_entry, enemy_wait := read_enemy_table()
	save_file(`./data/`, "enemy_entry", enemy_entry)
	save_file(`./data/`, "enemy_wait", enemy_wait)

	emu_port = "6503"
	send_command(memory_set(game_level_pointers, level_pointers))
	send_command(memory_set(game_level_data, level_data))

	send_command(memory_set(game_lantern_pointers, lantern_pointers))
	send_command(memory_set(game_lantern_data, lantern_data))

	send_command(memory_set(doors_next_level_addr, doors_next_level))
	send_command(memory_set(doors_player_pos_addr, doors_player_pos))

	send_command(memory_set(enemy_entry_addr, enemy_entry))
	send_command(memory_set(enemy_wait_addr, enemy_wait))

	cld := load_file(`./data/`, "code_loop_pointers", true)
	send_command(memory_set(game_code_init_pointers, cld))

	emu_port = "6502"

	screenshot()
}

func screenshot() {
	buffer_size = 150000
	display := send_command(display_get(), RESPONSE_DISPLAY_GET)
	raw_to_png("./img/", level, display)
	buffer_size = 50000
}

func detect_data_length(data, end_seq []byte) int {
	found := 0
	for i, _ := range data {
		j := i
		for _, s := range end_seq {
			if data[j] != s {
				found = 0
				break
			}
			found++
			j++
			if found == len(end_seq) {
				return i
			}
		}
	}
	return 0
}

/// s "map" 0 8c78 8fe7
/// s "color" 0 d878 dbff

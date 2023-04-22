package main

import (
	"os"
	"time"
)

var (
	//mem_map     = map[string]MemEntry{}
	map_path    = `./map/`
	level       = ""
	emu_port    = "6502" /// x64 -binarymonitoraddress ip4://127.0.0.1:6503
	buffer_size = 5000
	tcp_wait    = time.Duration(1) * time.Second
)

func main() {
	if len(os.Args) <= 2 {
		println("provide level number")
		return
	}
	if len(os.Args) != 3 {
		println("edit (e) or save (s) ?")
		return
	}
	level = os.Args[1]
	op := os.Args[2]
	switch op {
	case "e":
		load_editor(level)
	case "s":
		save_editor(level)
	}
}

func load_editor(level string) {
	decompress_editor_level(level)
}

func save_editor(level string) {
	compressed_level := compress_editor_level(level)
	save_file(map_path, "level"+level, compressed_level)

	lantern_records := read_lanterns_level(level)
	save_file(map_path, "lantern"+level, lantern_records)

	println("length: compressed_level|lantern_set:", len(compressed_level), len(lantern_records))

	level_pointers, level_data := link_byte(map_path, "level", game_level_data)
	save_file(map_path, "pointers_level", level_pointers)
	save_file(map_path, "data_level", level_data)

	lantern_pointers, lantern_data := link_word(map_path, "lantern", game_lantern_data)
	save_file(map_path, "pointers_lantern", lantern_pointers)
	save_file(map_path, "data_lantern", lantern_data)

	doors_next_level, doors_player_pos := read_doors_table()
	save_file(map_path, "doors_next_level", doors_next_level)
	save_file(map_path, "doors_player_pos", doors_player_pos)

	enemy_entry, enemy_wait := read_enemy_table()
	save_file(map_path, "enemy_entry", enemy_entry)
	save_file(map_path, "enemy_wait", enemy_wait)

	emu_port = "6503"
	memory_set(game_level_pointers, level_pointers)
	memory_set(game_level_data, level_data)

	memory_set(game_lantern_pointers, lantern_pointers)
	memory_set(game_lantern_data, lantern_data)

	memory_set(doors_next_level_addr, doors_next_level)
	memory_set(doors_player_pos_addr, doors_player_pos)

	memory_set(enemy_entry_addr, enemy_entry)
	memory_set(enemy_wait_addr, enemy_wait)

	cld := load_file(map_path, "code_loop_pointers", true)
	memory_set(game_code_init_pointers, cld)

	emu_port = "6502"

	/// screenshot()
}

/// s "map" 0 8c78 8fe7
/// s "color" 0 d878 dbff

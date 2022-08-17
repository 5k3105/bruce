package main

import (
	"fmt"
)

/// [level, enter, wait]
const (
	enemy_entry_addr = 0x4c9c
	enemy_wait_addr  = 0x4a8a
)

func read_enemy_table() ([]byte, []byte) { /// generate_doors_bin
	path := map_path
	data := load_file(path, "enemy.txt", false)
	lines, header := read_text_lines(data, "[")
	schema_seq := read_schema(header)
	dataset, num_rows := read_lines(lines, schema_seq)
	key_seq, field_seq := []string{"level"}, []string{"enter", "wait"}
	keymap := map_key(dataset, key_seq, field_seq, num_rows)
	return parse_enemy_map(keymap)
}

func parse_enemy_map(keymap map[string][]string) (enemy_entry, enemy_wait []byte) {
	for j := 1; j <= 20; j++ {
		key := fmt.Sprintf("%d", j)
		d, ok := keymap[key] /// field_seq: enter, wait
		if ok {
			s := "00"
			if d[0] == "f" {
				s = "ff"
			}
			enemy_entry = append(enemy_entry, htob(s))
			enemy_wait = append(enemy_wait, htob(d[1]))
		}
		if !ok {
			enemy_entry = append(enemy_entry, byte(0))
			enemy_wait = append(enemy_wait, byte(0))
		}
	}
	return enemy_entry, enemy_wait
}

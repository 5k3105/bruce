package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

/// [level(1-20), door(2-5), next-level, x, y]
const (
	doors_player_pos_addr = 0x4bfc
	doors_next_level_addr = 0x4b5c
)

/*
level-next-doors:
	doors 2-5
	20 levels each
	next level

sprite-pos:
	door 2-5
	20 levels each
	y pos
	x pos

*/

func read_doors_table() ([]byte, []byte) { /// generate_doors_bin
	path := map_path
	data := load_file(path, "doors.txt", false)
	lines, header := read_text_lines(data, "[")
	schema_seq := read_schema(header)
	dataset, num_rows := read_lines(lines, schema_seq)
	key_seq, field_seq := []string{"door", "level"}, []string{"next-level", "x", "y"}
	keymap := map_key(dataset, key_seq, field_seq, num_rows)
	return parse_door_map(keymap)
}

func read_schema(heading string) []string {
	heading = heading[1:]
	if heading[len(heading)-1:] == "]" {
		heading = heading[:len(heading)-1]
	}
	lines := strings.Split(heading, ",")

	var schema_seq []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			schema_seq = append(schema_seq, l)
		}
	}
	return schema_seq
}

func read_lines(lines []string, schema_seq []string) (map[string][]string, int) {
	dataset := map[string][]string{}
	for _, l := range lines {
		fields := strings.Fields(l)
		for i, f := range fields {
			f = strings.TrimSpace(f)
			if l != "" {
				key := schema_seq[i]
				dataset[key] = append(dataset[key], f)
			}
		}
	}
	return dataset, len(lines)
}

/// (door, level)  (next-level, x, y)
func map_key(dataset map[string][]string, key_seq, field_seq []string, num_rows int) map[string][]string {
	datamap := map[string][]string{}
	for i := 0; i < num_rows; i++ {
		keys := []string{}
		for _, key := range key_seq {
			keys = append(keys, dataset[key][i])
		}
		fields := []string{}
		for _, field := range field_seq {
			fields = append(fields, dataset[field][i])
		}
		ks := strings.Join(keys, ".")
		datamap[ks] = fields
	}
	return datamap
}

func parse_door_map(keymap map[string][]string) (next_level, sprite_pos []byte) {
	for i := 2; i <= 5; i++ {
		var sprite_posx, sprite_posy []byte
		for j := 1; j <= 20; j++ {
			key := fmt.Sprintf("%d.%d", i, j)
			d, ok := keymap[key] /// field_seq: next-level, x, y
			if ok {
				lv, _ := strconv.Atoi(d[0])
				next_level = append(next_level, uint8(lv-1))
				sprite_posx = append(sprite_posx, htob(d[1]))
				sprite_posy = append(sprite_posy, htob(d[2]))

			}
			if !ok {
				next_level = append(next_level, byte(0))
				sprite_posx = append(sprite_posx, byte(0))
				sprite_posy = append(sprite_posy, byte(0))
			}
		}
		sprite_pos = append(sprite_pos, sprite_posx...)
		sprite_pos = append(sprite_pos, sprite_posy...)
	}
	return next_level, sprite_pos
}

func htob(h string) byte {
	if len(h) == 1 {
		h = "0" + h
	}
	n, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return n[0]
}

func read_doors_bin() {
	path := map_path
	sprite_pos := load_file(path, "doors_player_pos", false)
	next_level := load_file(path, "doors_next_level", false)
	field_seq := []string{"level", "door", "next-level", "x", "y"}
	output_lines := generate_doors_text(next_level, sprite_pos, field_seq)
	doors_text := strings.Join(output_lines, "\n")
	save_file("", "doors.txt", []byte(doors_text))
}

func generate_doors_text(next_level, sprite_pos []byte, field_seq []string) []string {
	var output_lines []string
	header := strings.Join(field_seq, ", ")
	output_lines = append(output_lines, "["+header+"]")
	format := "%d\t%d\t%d\t%x\t%x"
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 19; j++ {
			line := fmt.Sprintf(format, j+1, i+2, next_level[(j+(i)*20)]+1, sprite_pos[(j+(i+i)*20)], sprite_pos[(j+(i+i+1)*20)])
			output_lines = append(output_lines, line)
		}
	}
	return output_lines
}

func read_text_lines(data []byte, mark string) ([]string, string) {
	contents := string(data)
	contents = strings.Replace(contents, string(byte(13)), "", -1)
	lf := string([]byte{10})
	lines := strings.Split(contents, lf)

	var lines2 []string
	var header string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			if string(l[0]) == mark {
				header = l
			} else {
				lines2 = append(lines2, l)
			}
		}
	}
	return lines2, header
}

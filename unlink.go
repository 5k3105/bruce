package main

import (
	"fmt"
)

/// split original dataset using pointer files

func split_linked_files() {
	path := `./original_dataset/`

	dataset := split_linked_data_byte("LevelPointers", "LevelData", game_level_slots)
	for i, data := range dataset {
		filename := fmt.Sprintf("level%d", i+1)
		save_file(path, filename, data)
	}

	dataset = split_linked_data_word("LanternPointers", "LanternData", game_level_slots)
	for i, data := range dataset {
		filename := fmt.Sprintf("lantern%d", i+1)
		save_file(path, filename, data)
	}

}

func split_linked_data_word(pointers_, data_ string, records int) [][]byte {
	path := `./original_data/`

	pointers := load_file(path, pointers_, false)
	data := load_file(path, data_, false)
	println("data len:", len(data))
	var cumu uint16
	dataset := [][]byte{}
	for i := 0; i < records*2-4; i += 2 {
		p1 := []byte{pointers[i], pointers[i+1]}
		p2 := []byte{pointers[i+2], pointers[i+3]}
		length := btoi16(p2) - btoi16(p1)
		cumu += length
		fmt.Printf("[%x][%x] [%d] [%d]\n", btoi16(p1), btoi16(p2), length, cumu)
		dataset = append(dataset, data[:length])
		data = data[length:]
	}
	dataset = append(dataset, data)
	println(len(dataset))
	return dataset
}

func split_linked_data_byte(pointers_, data_ string, records int) [][]byte {
	path := `./original_data/`

	pointers := load_file(path, pointers_, false)
	pointers_lo := pointers[:records]
	pointers_hi := pointers[records+1:]

	data := load_file(path, data_, false)
	dataset := [][]byte{}
	for i := 0; i < records-1; i++ {
		p1 := []byte{pointers_lo[i], pointers_hi[i]}
		p2 := []byte{pointers_lo[i+1], pointers_hi[i+1]}
		length := btoi16(p2) - btoi16(p1)
		fmt.Printf("[%x][%x] [%d]\n", btoi16(p1), btoi16(p2), length)
		dataset = append(dataset, data[:length])
		data = data[length:]
	}
	dataset = append(dataset, data)
	println(len(dataset))
	return dataset
}

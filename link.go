package main

import (
	"path/filepath"
	"strconv"
)

func link_word(path, filename string, start_address uint16) ([]byte, []byte) {
	files, err := filepath.Glob(path + filename + "*")
	if err != nil {
		panic(err)
	}

	if len(files) > game_level_slots {
		panic("len(files) > game_level_slots")
	}

	pointers := []byte{}
	dataset := []byte{}
	for i := 1; i < len(files); i++ {
		i_ := strconv.Itoa(i)
		fn := filename + i_
		println("f>  ", fn)
		data := load_file(path, fn, false)
		dataset = append(dataset, data...)

		pointers = append(pointers, i16tob(start_address)...)
		start_address += uint16(len(data))
	}

	return pointers, dataset
}

func link_byte(path, filename string, start_address uint16) ([]byte, []byte) {
	files, err := filepath.Glob(path + filename + "*")
	if err != nil {
		panic(err)
	}

	if len(files) > game_level_slots {
		panic("len(files) > game_level_slots")
	}

	pointers_lo := []byte{}
	pointers_hi := []byte{}

	dataset := []byte{}
	for i := 1; i < len(files); i++ {
		i_ := strconv.Itoa(i)
		fn := filename + i_
		data := load_file(path, fn, false)
		dataset = append(dataset, data...)

		bytes := i16tob(start_address)
		pointers_lo = append(pointers_lo, bytes[0])
		pointers_hi = append(pointers_hi, bytes[1])
		start_address += uint16(len(data))
	}

	pointers := []byte{}
	pointers = append(pointers, pointers_lo...)
	pad := 20 - len(pointers) + 1
	for i := 0; i < pad; i++ {
		pointers = append(pointers, 0xe0)
	}
	pointers = append(pointers, pointers_hi...)
	return pointers, dataset
}

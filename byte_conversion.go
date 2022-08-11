package main

/// https://gist.github.com/chiro-hiro/2674626cebbcb5a676355b7aaac4972d#file-golang-uint64-uint32-to-bytes-md

func i16tob(val uint16) []byte {
	r := make([]byte, 2)
	for i := uint16(0); i < 2; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

func btoi16(val []byte) uint16 {
	r := uint16(0)
	for i := uint16(0); i < 2; i++ {
		r |= uint16(val[i]) << (8 * i)
	}
	return r
}

func btoi32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func btoi64(val []byte) uint64 {
	r := uint64(0)
	for i := uint64(0); i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}

func b1tob(b bool) byte {
	byt := byte(0)
	if b {
		byt = 1
	}
	return byt
}

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func display_get() []byte {
	return create_command(0x84, []byte{1, 0})
}

func palette_get() []byte { /// no 0x91 response
	return create_command(0x91, []byte{1})
}

func memory_get(start_address uint16, length int) []byte {
	side_effects := byte(1)
	end_address := uint16(0)
	if end_address == 0 {
		end_address = start_address + uint16(length) - 1
	}
	memspace := byte(0)
	bank := uint16(0)

	b := []byte{side_effects}
	b = append(b, i16tob(start_address)...)
	b = append(b, i16tob(end_address)...)
	b = append(b, memspace)
	b = append(b, i16tob(bank)...)

	return create_command(0x01, b)
}

func memory_set(start_address uint16, content []byte) []byte {
	side_effects := byte(0)
	bank := uint16(0)
	memspace := byte(0)
	end_address := uint16(0)
	if end_address == 0 {
		end_address = start_address + uint16(len(content)) - 1
	}
	b := []byte{side_effects}
	b = append(b, i16tob(start_address)...)
	b = append(b, i16tob(end_address)...)
	b = append(b, memspace)
	b = append(b, i16tob(bank)...)
	b = append(b, content...)

	return create_command(0x02, b)
}

func jump(address uint16) {
	memspace := byte(0)
	array_count := uint16(0)
	item_size := byte(3)
	register_id := byte(3) /// PC

	b := []byte{memspace}
	b = append(b, i16tob(array_count)...)
	b = append(b, item_size)
	b = append(b, register_id)
	b = append(b, i16tob(address)...)

	send_command(create_command(0x32, b))
}

func send_command(cmd []byte, rc ...resp_code) []byte {
	c := &TCPClient{
		Host: "127.0.0.1",
		Port: emu_port,
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	conn.Write(cmd[:])

	data := []byte{}
	if rc != nil {
		time.Sleep(time.Duration(1) * time.Second)
		resp := step(conn)
		for _, r := range resp {
			if int(r.Req) == int(rc[0]) {
				data = parse_response_code(r.Body, rc[0])
			}
		}
	}

	conn.Close()
	return data
}

func parse_response_code(resp []byte, rc resp_code) []byte {
	data := []byte{}
	switch rc {
	case RESPONSE_MEM_GET:
		return response_mem_get(resp)
	case RESPONSE_DISPLAY_GET:
		return response_display_get(resp)
	case RESPONSE_PALETTE_GET:
		return response_palette_get(resp)
	}
	return data
}

func response_palette_get(resp []byte) []byte {
	num_items := int(btoi16(resp[0:2]))
	for i := 0; i <= num_items; i++ {
		fmt.Printf("rgb: ", resp[2+i:2+i+4])
		i += 4
	}
	println("-------------")
	return resp[2:]
}

func response_display_get(resp []byte) []byte {
	disp_buf_field_len := btoi32(resp[0:4])
	debug_width := btoi16(resp[4:6])
	debug_height := btoi16(resp[6:8])
	x_inner := btoi16(resp[8:10])
	y_inner := btoi16(resp[10:12])
	width_inner := btoi16(resp[12:14])
	height_inner := btoi16(resp[14:16])
	bits := int(resp[16])
	buffer_len := btoi32(resp[17:21])

	format := "disp_buf_field_len: [%d]\ndebug_width:        [%d]\ndebug_height:       [%d]\nx_inner:            [%d]\ny_inner:            [%d]\nwidth_inner:        [%d]\nheight_inner:       [%d]\nbits:               [%d]\nbuffer_len:         [%d]\n"
	fmt.Printf(format, disp_buf_field_len, debug_width, debug_height, x_inner, y_inner, width_inner, height_inner, bits, buffer_len)
	return resp[21:]
}

func response_mem_get(resp []byte) []byte {
	return resp[2:] ///	byte 0-1: The length of the memory segment.
}

func load_file_memory(path, filename string, start_address uint16) {
	content := load_file(path, filename, false)
	cmd := memory_set(start_address, content)
	println(filename, len(content))
	send_command(cmd)
}

func load_file(path, filename string, skip_load_addr bool) []byte {
	content, err := ioutil.ReadFile(path + filename)
	if err != nil {
		panic(err)
	}
	if skip_load_addr {
		content = content[2:]
	}
	return content
}

func save_file(path, filename string, data []byte) {
	err := os.WriteFile(path+filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func create_command(cmd byte, body []byte) []byte {
	c := []byte{0x02, 0x02}
	c = append(c, i32tob(uint32(len(body)))...)
	reqid := [4]byte{0xad, 0xde, 0x34, 0x12}
	c = append(c, reqid[:]...)
	c = append(c, cmd)
	c = append(c, body...)
	return c
}

func checkpoint_set(start_address, end_address uint16, stop_when_hit, enabled, temporary bool, cpu_operation string) []byte {
	memspace := byte(0)
	if end_address == 0 {
		end_address = start_address
	}

	b := []byte{}
	b = append(b, i16tob(start_address)...)
	b = append(b, i16tob(end_address)...)
	b = append(b, b1tob(stop_when_hit)) /// 0x01: true, 0x00: false
	b = append(b, b1tob(enabled))       /// 0x01: true, 0x00: false

	op := byte(0)
	switch cpu_operation {
	case "load":
		op = 1
	case "store":
		op = 2
	case "exec":
		op = 3
	}
	b = append(b, op)
	b = append(b, b1tob(temporary)) /// Deletes the checkpoint after it has been hit once. This is similar to "until" command, but it will not resume the emulator.
	b = append(b, memspace)         /// breaks? []byte{0xe2, 0xfc, 0xe3, 0xfc, 0x01, 0x01, 0x04, 0x01},

	return create_command(0x12, b)
}

func send_keyboardfeed(text string) []byte {
	b := []byte{}
	b = append(b, byte(len(text)))
	b = append(b, []byte(text)...)
	return create_command(0x72, b)
}

func send_continue(text string) []byte {
	b := []byte{0}
	return create_command(0xaa, b)
}

func step(conn *net.TCPConn) []*Response {
	rsp := []*Response{}

	reply := make([]byte, buffer_size)
	nbytes, err := conn.Read(reply)
	if err != nil {
		panic(err)
	}
	reply = reply[:nbytes]
	resp_len := 0
	for {
		fmt.Printf("[resp_len:nbytes] [%d]:[%d] ", resp_len, nbytes)
		resp := parse_response(reply[resp_len:nbytes]) /// add length check if buffer overflow
		rsp = append(rsp, resp)

		resp_error := decode_error(resp.Error)
		fmt.Printf("resp.Req: [%x][%v]\n", resp.Req, resp_error)

		resp_len += 12 + int(resp.Length)
		if resp_len >= nbytes {
			break
		}
	}

	return rsp
}

type Response struct {
	Stx    byte
	Api    byte
	Length uint32
	Req    byte
	Error  byte
	Reqid  uint32
	Body   []byte
}

func parse_response(resp []byte) *Response { // 12 + Len
	r := &Response{
		Stx:    resp[0],
		Api:    resp[1],
		Length: btoi32(resp[2:6]),
		Req:    resp[6],
		Error:  resp[7],
		Reqid:  btoi32(resp[8:12]),
	}
	r.Body = resp[12 : 12+r.Length] /// add length check if buffer overflow
	return r
}

func decode_error(resp byte) string {
	var err string
	switch resp {
	case 0x00:
		err = "OK, everything worked"
	case 0x01:
		err = "The object you are trying to get or set doesn't exist."
	case 0x02:
		err = "The memspace is invalid"
	case 0x80:
		err = "Command length is not correct for this command"
	case 0x81:
		err = "An invalid parameter value was present"
	case 0x82:
		err = "The API version is not understood by the server"
	case 0x83:
		err = "The command type is not understood by the server"
	case 0x8f:
		err = "The command had parameter values that passed basic checks, but a general failure occurred"
	}
	return err
}

type TCPClient struct {
	Host string
	Port string
}

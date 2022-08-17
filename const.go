package main

const (
	editor_compressed_level_addr = 0xc700
	editor_lantern_set_addr      = 0xc900 /// first byte is length + 1 ?
	editor_screen_addr           = 0x0478
	editor_screen_length         = 880
	editor_color_addr            = 0xd878
	editor_color_length          = 904
	editor_cset_addr             = 0x3800
	/// cset123 a000, a800 b000
)

const (
	game_level_slots    = 20
	game_level_pointers = 0x4b32 /// ?4 /// lo 20 hi 20
	game_level_data     = 0x68f9 /// - 817f LEN(6278)

	game_lantern_pointers = 0x4ce0 /// ?2 /// lo/hi word
	game_lantern_counts   = 0x4cce /// 20
	game_lantern_data     = 0x4d0a /// ?c /// - 4e94 $188 LEN(392)

	/// wait time bruce yamo ninja (per lvl)

	/// Color
	game_color_pointers = 0x4809 /// $480b lo/hi 20 start @ $47c8
	game_color_data     = 0x47c6 /// - 0x47f2 count:44 /// $47c8

	/// Raster (4 byte temp raster pointer table split)
	game_raster_pointers = 0x4831 /// $4833 lo/hi 20 start @ $47f9
	game_raster_data     = 0x47f3 /// - $4808 count:22 /// $47f9

	game_player_entry_positions = 0x4bae /// (XY) 20 + 20 X 1a Y c4
	game_enemy_entry_positions  = 0x4bd6 /// (XY) 20 + 20 X 9f Y 62

	game_code_init_pointers = 0x4e9a /// lo/hi 20 start @ $33f2
	game_code_init_data     = 0x33f2
	game_code_loop_pointers = 0x4ec2 /// lo/hi 20 start @ $3640
	game_code_loop_data     = 0x3640

	/*
		491E ; irq vec - main screen

	*/
)

type resp_code int

const (
	RESPONSE_INVALID              resp_code = 0x00
	RESPONSE_MEM_GET                        = 0x01
	RESPONSE_MEM_SET                        = 0x02
	RESPONSE_CHECKPOINT_INFO                = 0x11
	RESPONSE_CHECKPOINT_DELETE              = 0x13
	RESPONSE_CHECKPOINT_LIST                = 0x14
	RESPONSE_CHECKPOINT_TOGGLE              = 0x15
	RESPONSE_CONDITION_SET                  = 0x22
	RESPONSE_REGISTER_INFO                  = 0x31
	RESPONSE_DUMP                           = 0x41
	RESPONSE_UNDUMP                         = 0x42
	RESPONSE_RESOURCE_GET                   = 0x51
	RESPONSE_RESOURCE_SET                   = 0x52
	RESPONSE_JAM                            = 0x61
	RESPONSE_STOPPED                        = 0x62
	RESPONSE_RESUMED                        = 0x63
	RESPONSE_ADVANCE_INSTRUCTIONS           = 0x71
	RESPONSE_KEYBOARD_FEED                  = 0x72
	RESPONSE_EXECUTE_UNTIL_RETURN           = 0x73
	RESPONSE_PING                           = 0x81
	RESPONSE_BANKS_AVAILABLE                = 0x82
	RESPONSE_REGISTERS_AVAILABLE            = 0x83
	RESPONSE_DISPLAY_GET                    = 0x84
	RESPONSE_VICE_INFO                      = 0x85
	RESPONSE_PALETTE_GET                    = 0x91
	RESPONSE_JOYPORT_SET                    = 0xa2
	RESPONSE_USERPORT_SET                   = 0xb2
	RESPONSE_EXIT                           = 0xaa
	RESPONSE_QUIT                           = 0xbb
	RESPONSE_RESET                          = 0xcc
	RESPONSE_AUTOSTART                      = 0xdd
)

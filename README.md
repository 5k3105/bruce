```
Clone repo, 'go build'

Executable expects to see the /map folder
editor.vsf and bruce.vsf are needed for the emulator

Workflow:

 1. Launch 2 x64 (vice) emulators

	x64 -binarymonitoraddress ip4://127.0.0.1:6502 -keymap 1 editor.vsf

	x64 -binarymonitoraddress ip4://127.0.0.1:6503 -controlport1device 1 bruce_the_pro.vsf

 2. Use [level # 1-20] edit

	ex> ./bruce 1 e

	Loads the first level map into the editor. 
	Editor commands are listed below.

 3. Update ./map/doors.txt and enemy.txt if necessary (explained below).
	
 4. Use [level # 1-20] save 

	ex> ./bruce 1 s

	(s)ave the compiled level in the editor, link all other levels 
	and load them into the game to play test. Also reads doors and enemy
	text files and loads into game memory.

----------------------------------

Doors (doors.txt is in the /map folder)

 Pressing 9 in the editor will make the doors visible. In the game door #2-5
 will trigger a lookup into the 'doors' table:

 [level(1-20), door(2-5), next-level, x, y]

 level is current level
 door is the door # touched (door 1 has no effect)
 next-level is level to load next 
 x and y are sprite positions in hex where
	x increases from left to right
	y increases from top to bottom
	for position examples see doors.png and position_map.png

----------------------------------

Enemy (enemy.txt is in the /map folder)

 [level, enter, wait] 

 level is current level
 enter: if t, enemy can enter. if f, enemy cannot enter
 wait is time to wait before enemy appears (value in hex). 0 is no wait. 

----------------------------------

Editor command keys: (use positional keymap in vice)
	
	cursor : up, down, left right
	
	`   : cursor size (1/2/3 wide) 
	+/- : move chrset 
	a/s : black/blue color map (black: floor or obstacle, blue: pass thru or climb)
	z/x : fill/delete
	c/v : copy/paste
	123 : change multicolor background (screen area determined by raster row split / level #)
	F1  : show color map
	F5  : clear screen
	9   : show doors

----------------------------------

Limitations:

	This first version only saves (and links) level charcter map data, 
	color data and lantern data. Door links from level to level can be
	updated by text table and enemy entry and wait times by enemy.txt file.
	You can make many levels with lanterns and go from room to room.
	
Next features:

	Make background color rasters editable
	Set chrset per level
	Add lantern counts (will be used in code)
	Code edit and linking (currently all original level code is ignored)

Other notes:

	The disassembly 'brucelee_ida_dis.txt' was by me ~2008. 
	A more complete disassembly can be found by github.com/fa8ntomas
	for the Atari 8-bit and also has a game editor.

	If you would like to contribute to this version or have made some
	interesting levels please post in the Discussion area.

```	

```

Clone repo, 'go build'

Executable expects to see the /map /data and /img folders

Workflow:

1. Launch 2 x64 (vice) emulators

	x64 -binarymonitoraddress ip4://127.0.0.1:6502 -keymap 1 editor.vsf

	x64 -binarymonitoraddress ip4://127.0.0.1:6503 bruce.vsf

2. Use [level # 1-20] load

	ex> ./bruce 1 l

	(l)oads the first level map into the editor. 
	Editor commands are listed below.
	Make sure to hit F8 to compile the level before you use the save program.
	
3. Update ./map/doors.txt and enemy.txt if necessary (explained below).
	
4. Use [level # 1-20] save 

	ex> ./bruce 1 s

	(s)ave the compiled level in the editor, link all other levels 
	and load them into the game to play test. Also reads doors and enemy
	text files and loads into game memory. saves a screenshot in /img.

----------------------------------

Doors (doors.txt is in the /map folder)

Pressing 9 in the editor will make the doors visible. In the game door #2-5
will trigger a lookup into the 'doors' table:

[level(1-20), door(2-5), next-level, x, y]

level is current level
door is the door # touched
next-level is level to load next 
x and y are sprite positions in hex where
	x increases from left to right
	y increases from top to bottom
	for position examples see doors.png and position_map.png

(door 1 has no effect)

----------------------------------

Enemy (enemy.txt is in the /map folder)

level is current level
enter: if t, enemy can enter. if f, enemy cannot enter
wait is time to wait before enemy appears (value in hex). 0 is no wait. 

----------------------------------

Editor command keys: (use positional keymap in vice)

	+/-: move chrset 

	cursor keys: up, down, left right

	space: ?

	z/x: fill/delete

	123: change multicolor background (section determined by raster / level #)

	`: cursor size (1/2/3 wide) 

	F1: show color map

	F5: clear screen

	c/v: copy/paste

	0: move to char (seems broken...)

	9: show doors

	F8: compress level

	a/s: black/blue color map (black: floor, blue: pass thru or climb)

----------------------------------

Limitations:

	This first version only saves (and links) level charcter map data, 
	color data and lantern data. Door links from level to level can be
	updated by text table and enemy entry and wait times by enemy.txt file.
	You can make many levels with lanterns and go from room to room.
	
Next features:

	Make background color rasters editable.
	Add lantern counts (will be used in code)
	Code edit and linking (currently all original level code is ignored)
	
Bugs:

	There is a bug in the level compression code. Sometimes results in
	missing objects if not enough graphic blocks are on the level.
	ie If level is missing graphics, add more graphics. 
	

```	

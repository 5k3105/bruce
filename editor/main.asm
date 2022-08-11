    * = $C000

    .const lantern_records = $C900
    // .const LevRecAdr = $C700      // A3-B6?	NOT USED Level Data @ $C700

    .const zpRowCount = $A3
    .const zpRecStore1 = $A4      // First Record Storage
    .const zpRecStore2 = $A5      // Second Record Storage
    .const zpPageFlag = $A6       // Last Page Flag
    .const zpOutFlag = $A7        // RepOut/NorOut (0/1)
    .const zpEndFlag = $A8        // End Flag

    .const LevelPtr2 = $A9
    .const zpLevelLo2 = $A9
    .const zpLevelHi2 = $AA

    .const ScrnPtr = $AB
    .const zpScreenLo = $AB
    .const zpScreenHi = $AC

    .const ColrPtr = $AD
    .const zpColorLo = $AD
    .const zpColorHi = $AE

    .const LevelPtr = $AF
    .const zpLevelLo = $AF
    .const zpLevelHi = $B0

    //	raster compare: $E2 226 + 8 = $EA 234
    .const RasCmpL = $D012
    .const RasCmpH = $D011
    .const RasEnb = $D01A

    //	sprites
    .const Sp0X = $D000
    .const Sp0Y = $D001
    .const Sp1X = $D002
    .const Sp1Y = $D003
    .const SpMSB = $D010
    .const SpEnable = $D015

    .const SpPtr = $3700
    .const Sp0ptr = $7F8
    .const Sp1ptr = $7F9
    .const Sp2ptr = $BF8          // Color Map Page
    .const Sp3ptr = $BF9

    //	row static
    .const Row0 = $400            // $0770  //8C70			status row
    .const Row1 = $428            // $0798  //8C00
    .const Row2 = $450            // $07C0  //8C28

    .const Row0c = $D800          // $DB70  //D870			status row
    .const Row1c = $D828
    .const Row2c = $D850

    .const Xkey = $24
    .const OvTmp = $32

    // zp vectors - jump & video matrix
    .const VcLo = $68
    .const VcHi = $69
    .const VmLo = $6A
    .const VmHi = $6B             // video matrix
    .const JvLo = $6C
    .const JvHi = $6D             // jump vector

    // mode variables - temp store & state
    .const sinc = $6E
    .const SMst = $6E
    .const SMode = $6F            // state

    // mode zp vector & count
    .const SMvLo = $70
    .const SMvHi = $71
    .const SMct = $72             // selector mode

    // irq state, count, temp store
    .const IrqSwitch = $8B        // state
    .const ColorBuf = $8C         // ?
    .const IrqCnt = $8D           // unused?
    .const sprstx = $90
    .const sprsty = $91
    .const DSstore = $8E          // temp
    .const Paintcolor = $8F       // state

    .const ColorListIndx = $92
    .const RasLineIndx = $93
    .const CurrentColorIndx = $94
    .const MemPage = $95
    .const CopyBuf = $96          // Copy Buffer
    .const ChrBuf = $97           // Chr Buffer

    .const zp_lantern_index = $98 // Lantern Index
    .const zp_lantern_row = $99   // Lantern Row

    // Immediate Equates
    .const CmdKeyLen = $14
    .const PalletLen = $02

    // XY selector vars
    .const posx = $F7             // position
    .const posy = $F8
    .const Xinc = $F9             // 8
    .const Yinc = $FA             // 16
    .const Xbound = $FB           // 40
    .const Ybound = $FC           // 11

    .const Char = $FD
    .const Buff = $FE

    .const read_inp = $f142

    // initialize
    // Ras IRQ setup
    sei
    lda #$01
    sta RasEnb
    //LDA #$E1
    //STA RasCmpL
    lda #$1B
    sta RasCmpH
    //STA IrqSwitch
    lda #<RasIrq
    sta $0314
    lda #>RasIrq
    sta $0315

    // disable hw interrupt
    // keyboard interrupt keeps running
    lda #$7F
    sta $DC0D
    sta $DD0D
    cli

clear_sprite_area:

    ldx #$00
    lda #$00

    {
        sta SpPtr, x
        inx
        cpx #$40
        bne -
    }

    lda #$DC
    sta Sp0ptr
    sta Sp1ptr
    sta Sp2ptr
    sta Sp3ptr

    lda #$03
    sta SpEnable

    // Sprite Positions
    lda #$18                      // Left
    sta sprstx
    sta Sp0X                      // $D000
    lda #$4A                      // 3A	//Top
    sta sprsty
    sta Sp0Y                      // $D001

    lda #$B0
    sta Sp1X                      // $D002
    lda #$3A                      // EA
    sta Sp1Y                      // $D003

    lda #$00
    sta SpMSB                     // $D010

    // mem-ctrl: chrset3800,screen400
    //LDA #$1F
    //STA $D018
    //LDA $D016
    //ORA #%10000
    //STA $D016	//multi-color
    // xy char positions
    lda #$00
    sta posx                      // 00-27 (40)
    sta posy                      // 00-0A (11)
    lda #$08                      // 8
    sta Xinc
    lda #$10                      // 16
    sta Yinc
    lda #$28                      // 40
    sta Xbound
    lda #$0B                      //11 = 22
    sta Ybound

    // initialize cursor mode (123)
    lda #$FF
    sta SMode
    jsr Mode

    // key auto-repeat
    lda #$80
    sta $028A

    // skip door displays
    lda #$05
    sta Char
    sta Buff
    jsr ClearScreen

color_screen:

    lda #$08
    sta Paintcolor
    lda #$00
    sta $D020                     // border color black
    ldx #$28

    {
        dex
        sta Row0c, x
        bne -
    }

color_pallet:

    ldx #PalletLen
    ldy #$00

    {
        lda #$A0
        sta Row0 + 21, y
        lda ColorPallet, y
        sec
        sbc #$08
        sta Row0c + 21, y
        iny
        dex
        bne -
    }

    lda #$1F
    sta MemPage
    jsr ColorMapSetup

    // input
GETINPUT:
    {
    L1: jsr read_inp
        beq L1

        ldx #$00
    L2: cmp KeyM, x
        beq L3
        inx
        cpx #CmdKeyLen                // length + 1
        bne L2

        ldx #$00
    L5: cmp ColorKey, x
        beq L6
        inx
        cpx #PalletLen                // length + 1
        bne L5
        beq L1

    L3: txa
        asl
        tax
        lda JmpM, x                   // init vector - reusable JvLo/Hi
        sta JvLo
        pha
        inx
        lda JmpM, x
        sta JvHi
        pha
        ldx #$00
        jsr L4
        pla
        pla
        jsr DrawStatus
        clc
        bcc L1

    L4: jmp (JvLo)

        // StoreColor:
    L6: lda ColorPallet, x
        sta Paintcolor
        ldx #$28

        {
            dex
            sta Row1c, x
            sta Row2c, x
            bne -
        }

        jsr GetV
        jsr DrCol
        clc
        lda VmLo
        adc #$28

        {
            bcc +
            inc VmHi
            inc VcHi
        }

        sta VmLo
        jsr DrCol

        clc
        bcc L1
    }
    // status bar
DrawStatus:

    ldy #$00
    lda #$18                      // 'X'
    sta Row0, y
    lda posx
    sta DSstore
    iny
    jsr DSwrite
    iny
    iny
    lda #$19                      // 'Y'
    sta Row0, y
    iny
    lda posy
    sta DSstore
    jsr DSwrite

    iny
    iny                           // chrset char
    lda Row1 + 19
    sta DSstore
    jsr DSwrite

    jsr GetV                      // char under cursor
    ldy #$00
    lda (VmLo), y
    sta DSstore
    ldy #$0E
    jsr DSwrite

    //LDY #$22  		//1E - colors
    lda #$A0
    sta DSstore
    //LDX #$00
    jsr DSwcolor
    rts

DSwcolor:

    ldy #$26
    lda DSstore
    sta Row0, y
    lda Paintcolor
    sta Row0c, y
    rts

DSwrite:

    lda #$24                      // '$'
    sta Row0, y
    lda DSstore
    and #$F0
    clc
    ror
    ror
    ror
    ror
    tax
    lda HexCodes, x
    iny
    sta Row0, y
    lda DSstore
    and #$0F
    tax
    lda HexCodes, x
    iny
    sta Row0, y
    rts

    // colors
One:

    ldx CurrentColorIndx
    inc ColorList, x              // $D021
    rts
Two:

    ldx CurrentColorIndx
    inc ColorList + 1, x          // $D022
    rts
Three:

    ldx CurrentColorIndx
    inc ColorList + 2, x          // $d023
    rts

    // fill/space
    // get screen vector using table
    // input posy,VcTbl,VlTbl,posx
GetV:
    {
        ldx posy                      //row
        lda VhTbl, x
        sta VmHi                      //video
        lda VcTbl, x
        sta VcHi                      //color
        lda VlTbl, x
        clc
        adc posx                      //column
        bcc L1
        inc VmHi
        inc VcHi
    L1: sta VmLo
        rts                           //returns VmHi,VmLo,VcHi,posyx
    }
Fill:

    lda Row1 + 19                 //selector character
    sta ChrBuf
    jsr WriteChr
    rts

WriteChr:
    {
        jsr GetV
        lda ChrBuf                    // selector character
        jsr DrVid                     // write

        lda VmHi
        pha
        jsr DrCol

        pla
        sta VmHi
        lda VmLo
        adc #$28
        bcc L3
        inc VmHi
        inc VcHi
    L3: sta VmLo

        lda ChrBuf
        clc
        adc #$80
        jsr DrVid
    }

DrCol:
    {
        lda VcHi
        sta VmHi
        ldx SMode
        inx
        lda Paintcolor
        ldy #$00
    L1: sta (VmLo), y
        iny
        dex
        bne L1
        rts
    }

DrVid:
    {
        ldx SMode
        inx
        ldy #$00
        clc
    L2: sta (VmLo), y
        adc #$01
        iny
        dex
        bne L2
        rts
    }

Space:
    {
        jsr GetV
        lda #$00
        jsr DrSp
        lda VmLo
        adc #$28
        bcc L1
        inc VmHi
    L1: sta VmLo
        lda #$80
    }

DrSp:
    {
        ldx SMode
        inx
        ldy #$00
        clc
    L2: sta (VmLo), y
        iny
        dex
        bne L2
        rts
    }
    // plus/minus
Plus:
    {
        inc Buff
        bpl L1
        lda #$01                      // #$06	//#$05
        sta Buff
    L1: lda Buff
        sta Char
        jsr Lc1
        rts
    }

Minus:
    {
        dec Buff
        lda Buff
        cmp #$00                      // #$05
        bne L1
        lda #$7F
        sta Buff
    L1: lda Buff
        sta Char
        jsr Lc1
        rts
    }
pm_end:

    lda #$00                      // #$05
    sta Char

Lc1:

    inc Char
    lda Char                      // 06 START
    bmi pm_end                    // 80
    sta Row1, x
    clc
    adc #$80
    sta Row2, x
    inx
    cpx #$28
    bne Lc1
    rts

    // sprite mode
Mode:

    lda #<Sptbl
    sta SMvLo
    lda #>Sptbl
    sta SMvHi

    inc SMode
    ldx SMode
    cpx #$03
    bne Modetop
    lda #$00
    sta SMode
    ldx #$00

Modetop:

    lda SMvLo
    clc
    adc SMtbl, x
    sta SMvLo
    bcc SM
    inc SMvHi

SM: {
        sei
        ldx #$00
    L1: lda SMcnt, x
        pha
        inx
        cpx #$03
        bne L1                        // 01,2D,01
        stx SMst                      // 3 records

        ldx #$00
    L5: pla
        sta SMct
    L4: ldy #$00
    L2: lda (SMvLo), y
        sta SpPtr, x
        inx
        iny
        cpy #$03                      // 3 bytes
        bne L2
        dec SMct
        bne L4

        lda SMvLo
        clc
        adc #$03                      // increment record
        sta SMvLo
        bcc L3
        inc SMvHi

    L3: dec SMst
        bne L5
        cli
        rts
    }

    // data
SMcnt:

    .byte $01, $0E, $01

SMtbl:

    .byte $00, $09, $12

Sptbl:

    .byte $FF, $00, $00
    .byte $81, $00, $00
    .byte $FF, $00, $00

    .byte $FF, $FF, $00
    .byte $80, $01, $00
    .byte $FF, $FF, $00

    .byte $FF, $FF, $FF
    .byte $80, $00, $01
    .byte $FF, $FF, $FF

VcTbl:

    .byte $D8, $D8
    .byte $D9, $D9, $D9
    .byte $DA, $DA, $DA, $DA
    .byte $DB, $DB

VhTbl:

    .byte $04, $04
    .byte $05, $05, $05
    .byte $06, $06, $06, $06
    .byte $07, $07

VlTbl:

    .byte $78, $C8
    .byte $18, $68, $B8
    .byte $08, $58, $A8, $F8
    .byte $48, $98

    //CmdKeyLen	= $14 - Create Lant Recs key now unnecessary (F7)
KeyM:
    //    +	   -    up   dn   lf   rt   sp   z    1    2    3    `    F1   F5
    .byte $2D, $2B, $91, $11, $9D, $1D, $58, $5A, $31, $32, $33, $5F, $85, $87
    //    c    v    0    9    F7   F8
    .byte $43, $56, $30, $39, $88, $8C

JmpM:

    .word Plus, Minus, up, dn, lf, rt, Space, Fill, One, Two, Three, Mode, ColorMap, ClearScreen
    .word Copy, Paste, MoveToChr, ShowDoors, create_lantern_records, Compress

HexCodes:

    .byte $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
    .byte $01, $02, $03, $04, $05, $06

RowRas:

    .byte $02, $04, $0A

ColorList:

    .byte $FB, $F1, $FE, $FB, $F1, $FE, $FB, $F2, $F1, $FC, $F2, $F1// 4x 3(# rasters)

RasLine:

    .byte $3A, $49, $78, $97, $FD, $78, $97, $FD

    // PalletLen 	= $02
ColorPallet:

    .byte $08, $0E                // black=block, blue=climb

    //ColorKeyLen = PalletLen
ColorKey:
    //     A   S   D   F   G   H   J   K
    .byte $41, $53, $44, $46, $47, $48, $4A, $4B

lantern_lookup_table:

    .byte $13, $14, $31, $9E, $9F, $92

    // IRQ
RasIrq:
    {
        lda #$01
        sta $D01A                     // raster cmp
        sta $D019                     // raster ack

        lda $D012
        cmp #$FD
        bne L1

        jsr $EA87                     // read keyboard

        lda $D016
        and #%11101111
        sta $D016                     //multi-color off

        lda #$15                      // regular chrset
        sta $D018
        lda #$FC                      // med grey
        sta $D021                     // status line
        bne L9

        //L1	TAX : INX : INX : TXA : INX
        //L10	CPX $D012
        //		BEQ L12
        //		CMP $D012
        //		BNE L10
        //		LDX #$0A
        //L11	DEX
        //		BNE L11
        //	LDA RasLineIndx
        //	BNE .L5
    L1: lda RasLineIndx
        bne L5

        lda MemPage                   // txt pg & chrset
        sta $D018

        lda $D016
        ora #%10000
        sta $D016                     //multi-color on

        ldx CurrentColorIndx
        bne L7

    L5: ldx ColorListIndx
    L7: lda ColorList, x
        sta $D021
        lda ColorList + 1, x
        sta $D022
        lda ColorList + 2, x
        sta $D023

        lda ColorListIndx             // next color record
        clc
        adc #$03
        sta ColorListIndx

        inc RasLineIndx               // next raster record
    L9: ldx RasLineIndx
        lda RasLine, x

    L2: sta $D012
        cmp #$FD                      // end of screen
        beq L3

    L4: pla
        tay
        pla
        tax
        pla
        rti

    L3: lda #$00
        sta RasLineIndx
        sta ColorListIndx
        jmp L4
    }

    // direction (l/r/u/d)
direction_move:

lf: {
        ldx posx
        dex
        bpl L1                        // >#$80?
        ldx #$27
    L1: stx posx
        jsr sprposinitx
        rts
    }

rt: {
        ldx posx
        inx
        cpx Xbound
        bcc L2                        // <#$28?
        ldx #$00
    L2: stx posx
        jsr sprposinitx
        rts
    }

up: {
        ldx posy
        dex
        bpl L3
        ldx #$0A
    L3: stx posy
        jsr sprposinity
        rts
    }

dn: {
        ldx posy
        inx
        cpx Ybound
        bcc L4
        ldx #$00
    L4: stx posy
        jsr sprposinity
        rts
    }

sprposinitx:
    {
        ldy #$00
        sty SpMSB
        txa
        asl
        asl
        asl
        bcc L5
        inc SpMSB
        adc #$17
        sta Sp0X
        rts
    L5: adc sprstx
        bcc L6
        inc SpMSB
    L6: sta Sp0X
        rts
    }

sprposinity:

    txa
    asl
    asl
    asl
    asl
    adc sprsty
    sta Sp0Y
    jsr ChrColorBar
    rts

ChrColorBar:
    {
        ldy #$FF
        ldx #$00
    L1: inx
        inx
        inx                           //set chrset bar colors
        iny
        lda RowRas, y
        cmp posy
        bcc L1
        stx CurrentColorIndx

        ldy #$00                      //save to chrset color bar buffer
    L2: lda #$A0                      //& update status
        sta Row0 + 34, y
        lda ColorList, x
        sta ColorList, y
        sta Row0c + 34, y
        inx
        iny
        cpy #$03
        bne L2

        rts
    }

    // F1
ColorMap:
    {
        lda MemPage
        cmp #$1F
        beq L1
        lda #$1F
        sta MemPage
        rts
    L1: lda #$2F
        sta MemPage
        rts
    }

    // Set up Color Map Page
ColorMapSetup:
    {
        ldx #$11

    L2: lda VhTbl, x
        adc #$04
        sta VmHi
        lda VlTbl, x
        sta VmLo

        ldy #$51                      // 81

    L1: dey
        lda #$06
        sta (VmLo), y
        tya
        bne L1

        dex
        bpl L2
        rts
    }

    // F5
ClearScreen:
    {
        ldx #$0A
    L2: lda VhTbl, x
        sta VmHi
        lda VlTbl, x
        sta VmLo

        ldy #$51                      // 81
        lda #$00
    L1: dey
        sta (VmLo), y
        bne L1
        dex
        bpl L2

        ldx #$0A
    L3: lda VcTbl, x
        sta VmHi
        lda VlTbl, x
        sta VmLo

        ldy #$51                      // 81
        lda #$0E
    L4: dey
        sta (VmLo), y
        bne L4
        dex
        bpl L3
        rts
    }

    // Copy/Paste (C/V)
Copy:

    jsr GetV                      //char under cursor
    ldy #$00
    lda (VmLo), y
    sta CopyBuf
    rts

Paste:

    lda CopyBuf
    sta ChrBuf
    jmp WriteChr

    // 0
MoveToChr:

    jsr GetV                      //char under cursor
    ldy #$00
    lda (VmLo), y
    //		CLC
    //		SBC #$10
    //		STA Buff
    //		INC Buff
    sta Char
    jsr Lc1
    //		JSR Minus
    rts

    // Import Charcters 1-5
    // 9
ShowDoors:
    {
        sei
        lda $01
        and #$FB
        sta $01

        ldx #$FF
    L1: inx
        lda $D188, x
        sta $3808, x
        sta $3C08, x
        cpx #$27
        bne L1

        lda $01
        ora #$04
        sta $01
        cli
        rts
    }

    // moved to Compressor .. F7 = #$88
create_lantern_records:
    {
        lda #$00
        sta zp_lantern_index
        lda #$ff
        sta zp_lantern_row

    L3: inc zp_lantern_row
        //inc zp_lantern_row
        lda zp_lantern_row
        cmp #$0c                      // $0A is last row
        beq L4                        // rts

        ldx zp_lantern_row
        lda VhTbl, x
        sta VmHi
        lda VlTbl, x
        sta VmLo

        ldy #$ff                      // column
    L2: cpy #$27                      // (VmLo),y
        beq L3                        // next row
        ldx #$ff                      // lantern_lookup_table,x
        iny                           // next column
        lda (VmLo), y
    L1: inx
        cmp lantern_lookup_table, x
        beq L5
        cpx #$02                      // 0->2 = 3
        bne L1
        beq L2

    L5: inx
        inx
        inx
        lda lantern_lookup_table, x
        pha
        lda #$01
        inc zp_lantern_index
        ldx zp_lantern_index
        sta lantern_records, x        // ON/OFF

        pla
        inc zp_lantern_index
        ldx zp_lantern_index
        sta lantern_records, x        // image stored

        tya
        inc zp_lantern_index
        ldx zp_lantern_index
        sta lantern_records, x        // column

        lda zp_lantern_row
        adc #$01
        inc zp_lantern_index
        ldx zp_lantern_index
        sta lantern_records, x        // row + 2

        clc
        bcc L2

    L4: lda #$ff
        inc zp_lantern_index
        ldx zp_lantern_index
        sta lantern_records, x
        lda zp_lantern_index
        sta lantern_records           // length 1st byte
        rts
    }
    // compressor
    // jsr Create Lant Recs after finished compiling level
    // F8 = #$8C
Compress:                         // SEI

    lda #$27
    sta zpRowCount
    lda #$04
    sta zpScreenHi
    lda #$78
    sta zpScreenLo                // $0478
    lda #$D8
    sta zpColorHi
    lda #$78
    sta zpColorLo                 // $D878
    lda #$C7
    sta zpLevelHi

    ldy #$00
    sty zpPageFlag
    sty zpOutFlag
    sty zpEndFlag
    sty zpLevelLo

Top:

    lda zpEndFlag
    bne End
    ldx #$00
    jsr Fetch
    sta zpRecStore1
    inx
    jsr IncScrnPtr
    beq Repeater

    inc zpOutFlag
    lda zpLevelLo
    sta zpLevelLo2
    lda zpLevelHi
    sta zpLevelHi2
    jsr IncLevPtr

Nor:

    lda zpRecStore1
    sta (LevelPtr), y
    jsr IncLevPtr
    inx
    bmi NorOut
    lda zpRecStore2
    sta zpRecStore1
    jsr IncScrnPtr
    bne Nor
NorOut:

    dex
    txa
    ora #$80
    sta (LevelPtr2), y
    dec zpOutFlag
    lda zpEndFlag
    bne End
    ldx #$01

Repeater:

    inx
    bmi RepOut
    jsr IncScrnPtr
    beq Repeater
RepOut:

    txa
    sta (LevelPtr), y
    jsr IncLevPtr
    lda zpRecStore1
    sta (LevelPtr), y
    jsr IncLevPtr
    jmp Top

End:

    lda #$80
    sta (LevelPtr), y
    jsr create_lantern_records
    rts
    // first byte of lantern_records = length
    // zpLevelLo/Hi = End Adr
IncScrnPtr:

    dec zpRowCount
    bpl RowOut
    lda #$27
    sta zpRowCount
    lda zpScreenLo
    clc
    adc #$28
    sta zpScreenLo
    sta zpColorLo
    bcc RowOut
    inc zpScreenHi
    inc zpColorHi

RowOut:

    lda zpPageFlag
    bne LastPage
    inc zpScreenLo
    inc zpColorLo
    bne Fetch
    inc zpScreenHi
    inc zpColorHi
    lda zpScreenHi
    cmp #$07
    bne Fetch
    inc zpPageFlag

Fetch:

    lda (ColrPtr), y
    and #$0F
    cmp #$08
    bne NoColr
    lda (ScrnPtr), y
    ora #$80
    sta zpRecStore2
    cmp zpRecStore1
    rts

NoColr:

    lda (ScrnPtr), y
    sta zpRecStore2
    cmp zpRecStore1
    rts

Exit:

    inc zpEndFlag
    pla
    pla
    lda zpOutFlag
    beq RepOut
    bne NorOut

LastPage:

    inc zpScreenLo
    inc zpColorLo
    lda zpScreenLo
    cmp #$E8
    beq Exit
    bne Fetch

IncLevPtr:

    inc zpLevelLo
    bne LevPtrOut
    inc zpLevelHi
LevPtrOut:

    rts

    // note: to find color ram in vice, search for 0 hi nyble instead of f ie.. 02 04 06 08 01 03 05 07 09

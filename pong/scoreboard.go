package main

import "github.com/nsf/termbox-go"

type Scoreboard struct {
	coords Pos
}

const (
	RIGHT_COL int = iota
	LEFT_COl
	FIRST_ROW
	SECOND_ROW
)

func (gm *GameState) drawScoreboard() {
	gm.drawAscii(P, Ascii_Display{FIRST_ROW, LEFT_COl}, PLAYER_ONE)
	gm.drawAscii(one, Ascii_Display{FIRST_ROW, RIGHT_COL}, PLAYER_ONE)
	gm.drawAscii(zero, Ascii_Display{SECOND_ROW, LEFT_COl}, PLAYER_ONE)
	gm.drawAscii(seven, Ascii_Display{SECOND_ROW, RIGHT_COL}, PLAYER_ONE)

	gm.drawAscii(I, Ascii_Display{FIRST_ROW, LEFT_COl}, PLAYER_TWO)
	gm.drawAscii(A, Ascii_Display{FIRST_ROW, RIGHT_COL}, PLAYER_TWO)
	gm.drawAscii(three, Ascii_Display{SECOND_ROW, LEFT_COl}, PLAYER_TWO)
	gm.drawAscii(one, Ascii_Display{SECOND_ROW, RIGHT_COL}, PLAYER_TWO)
}
func (gm *GameState) drawAsciiDivider() {

}

func (gm *GameState) drawAscii(ascii Ascii, display Ascii_Display, id int) {
	x_offset := asciiXOffset(display[1])
	y_offset := asciiYOffset(display[0])

	if id == PLAYER_TWO {
		for y, row := range mirrorAscii(ascii) {
			for x := range row {
				pos := Pos{x: x, y: y}
				char := row[x]
				gm.drawAsciiCell(pos, char, [2]int{y_offset, x_offset}, id)
			}
		}
		return
	}
	for y, row := range ascii {
		for x := range row {
			pos := Pos{x: x, y: y}
			char := row[x]
			gm.drawAsciiCell(pos, char, [2]int{y_offset, x_offset}, id)
		}
	}
}
func (gm *GameState) drawAsciiCell(pos Pos, char rune, offset [2]int, id int) {
	termbox.SetCell(gm.asciiXPos(pos.x, offset[1], id), gm.asciiYPos(pos.y, offset[0]), char, termbox.ColorGreen, termbox.ColorDefault)
}
func (gm *GameState) asciiXPos(x int, offset int, id int) int {
	if id == PLAYER_TWO {
		return gm.terminal.width - gm.scorboard.coords.x - offset - x
	}
	return gm.scorboard.coords.x + x + offset
}
func (gm *GameState) asciiYPos(y int, offset int) int {
	return gm.scorboard.coords.y + y + offset
}
func asciiXOffset(col int) int {
	if col == LEFT_COl {
		return 0
	}
	return Ascii_Width + Gap
}
func asciiYOffset(row int) int {
	if row == FIRST_ROW {
		return 0
	}
	return Ascii_Height + Gap
}
func mirrorAscii(ascii Ascii) Ascii {
	rv := Ascii{}

	for y, row := range ascii {
		for x := range row {
			idx := len(row) - 1 - x
			rv[y][idx] = row[x]
		}
	}
	return rv
}

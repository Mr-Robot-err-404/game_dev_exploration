package main

import "github.com/nsf/termbox-go"

type Scoreboard struct {
	padding Pos
}

func (gm *GameState) drawAscii(ascii Ascii) {
	for y, row := range ascii {
		for x := range row {
			pos := Pos{x: x, y: y}
			char := row[x]
			gm.drawAsciiCell(pos, char)
		}
	}
}

func (gm *GameState) drawAsciiCell(pos Pos, char rune) {
	termbox.SetCell(gm.asciiXPos(pos.x), gm.asciiYPos(pos.y), char, termbox.ColorGreen, termbox.ColorDefault)
}

func (gm *GameState) asciiXPos(x int) int {
	return gm.padding.x - gm.scorboard.padding.x + x
}

func (gm *GameState) asciiYPos(y int) int {
	return gm.padding.y + gm.scorboard.padding.y + y
}

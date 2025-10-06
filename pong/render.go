package main

import "github.com/nsf/termbox-go"

func (gm *GameState) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	gm.createBoard()
	gm.drawPlayer(gm.player.position, PLAYER_ONE)
	gm.drawPlayer(gm.opponent.position, PLAYER_TWO)
	gm.drawBall()

	termbox.Flush()
}
func (gm *GameState) createBoard() {
	char := '█'
	for i := 0; i <= Board.width; i++ {
		pos := Pos{x: i, y: 0}
		gm.drawCell(pos, char)

		pos.y = gm.board.height
		gm.drawCell(pos, char)
	}
	for i := 0; i <= Board.height; i++ {
		pos := Pos{x: 0, y: i}
		gm.drawCell(pos, char)

		pos.x = gm.board.width
		gm.drawCell(pos, char)
	}
}

func (gm *GameState) drawPlayer(pos Pos, id int) {
	y := pos.y / 2
	x := pos.x / 2
	mapped := Pos{x: x, y: y}

	switch gm.orientation {
	case STD:
		gm.drawStandardPlayer(mapped, id, pos.y)
	case ALT:
		gm.drawAltPlayer(mapped, id, pos.x)
	}
}

func (gm *GameState) drawAltPlayer(pos Pos, id int, x int) {
	if isEven(x) {
		if id == PLAYER_TWO {
			gm.drawCell(pos, '▀')
			gm.drawCell(Pos{x: pos.x + 1, y: pos.y}, '▀')
			gm.drawCell(Pos{x: pos.x + 2, y: pos.y}, '▀')
			gm.drawCell(Pos{x: pos.x + 3, y: pos.y}, '▀')
			return
		}
		gm.drawCell(pos, '▄')
		gm.drawCell(Pos{x: pos.x + 1, y: pos.y}, '▄')
		gm.drawCell(Pos{x: pos.x + 2, y: pos.y}, '▄')
		gm.drawCell(Pos{x: pos.x + 3, y: pos.y}, '▄')
		return
	}
	if id == PLAYER_TWO {
		gm.drawCell(pos, '▝')
		gm.drawCell(Pos{x: pos.x + 1, y: pos.y}, '▀')
		gm.drawCell(Pos{x: pos.x + 2, y: pos.y}, '▀')
		gm.drawCell(Pos{x: pos.x + 3, y: pos.y}, '▀')
		gm.drawCell(Pos{x: pos.x + 4, y: pos.y}, '▘')
		return
	}
	gm.drawCell(pos, '▗')
	gm.drawCell(Pos{x: pos.x + 1, y: pos.y}, '▄')
	gm.drawCell(Pos{x: pos.x + 2, y: pos.y}, '▄')
	gm.drawCell(Pos{x: pos.x + 3, y: pos.y}, '▄')
	gm.drawCell(Pos{x: pos.x + 4, y: pos.y}, '▖')
}

func (gm *GameState) drawStandardPlayer(pos Pos, id int, y int) {
	if isEven(y) {
		if id == PLAYER_ONE {
			gm.drawCell(pos, '▐')
			gm.drawCell(Pos{x: pos.x, y: pos.y + 1}, '▐')
			return
		}
		gm.drawCell(pos, '▌')
		gm.drawCell(Pos{x: pos.x, y: pos.y + 1}, '▌')
		return
	}
	if id == PLAYER_ONE {
		gm.drawCell(pos, '▗')
		gm.drawCell(Pos{x: pos.x, y: pos.y + 1}, '▐')
		gm.drawCell(Pos{x: pos.x, y: pos.y + 2}, '▝')
		return
	}
	gm.drawCell(pos, '▖')
	gm.drawCell(Pos{x: pos.x, y: pos.y + 1}, '▌')
	gm.drawCell(Pos{x: pos.x, y: pos.y + 2}, '▘')
}

func (gm *GameState) drawBall() {
	pos := gm.ball.position
	q := mapQuadrant(pos.x, pos.y)

	x, y := pos.x/2, pos.y/2
	mapped := Pos{x: x, y: y}

	switch q {
	case TOP_LEFT:
		gm.drawCell(mapped, '▘')
	case TOP_RIGHT:
		gm.drawCell(mapped, '▝')
	case BOTTOM_LEFT:
		gm.drawCell(mapped, '▖')
	case BOTTOM_RIGHT:
		gm.drawCell(mapped, '▗')
	}
}
func (gm *GameState) drawCell(pos Pos, char rune) {
	termbox.SetCell(gm.xPos(pos.x), gm.yPos(pos.y), char, termbox.ColorGreen, termbox.ColorDefault)
}
func (gm *GameState) xPos(x int) int {
	return gm.padding.x + x
}
func (gm *GameState) yPos(y int) int {
	return gm.padding.y + y
}

func mapQuadrant(x int, y int) int {
	if isEven(x) && isEven(y) {
		return TOP_LEFT
	}
	if isEven(x) {
		return BOTTOM_LEFT
	}
	if isEven(y) {
		return TOP_RIGHT
	}
	return BOTTOM_RIGHT
}

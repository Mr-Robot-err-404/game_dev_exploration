package main

// TODO: draw players
// move player
// draw ball
// move ball
// collision with wall
// collision with player
// ai logic

import (
	"time"

	"github.com/nsf/termbox-go"
)

type Pos struct {
	x int
	y int
}
type Grid struct {
	width  int
	height int
}
type GameState struct {
	board    Grid
	terminal Grid
	padding  Pos
	player   Pos
	opponent Pos
}
type Draw struct {
	pos  Pos
	char rune
}

func (gm *GameState) drawPlayer(pos Pos) {
	gm.drawCell(pos, '█')
	gm.drawCell(Pos{x: pos.x, y: pos.y + 1}, '█')
}

func (gm *GameState) drawCell(pos Pos, char rune) {
	termbox.SetCell(gm.xPos(pos.x), gm.yPos(pos.y), char, termbox.ColorGreen, termbox.ColorBlack)
}
func (gm *GameState) xPos(x int) int {
	return gm.padding.x + x
}
func (gm *GameState) yPos(y int) int {
	return gm.padding.y + y
}

var Board = Grid{width: 60, height: 18}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	terminal := setup()
	padding := getPadding(terminal)
	game := GameState{
		board:    Board,
		terminal: terminal,
		padding:  padding,
		player:   Pos{x: 2, y: Board.height / 3},
		opponent: Pos{x: Board.width - 2, y: Board.height / 2},
	}
	game.createBoard()
	game.drawPlayer(game.player)
	game.drawPlayer(game.opponent)
	termbox.Flush()

	time.Sleep(time.Second * 5)
}

func setup() Grid {
	width, height := termbox.Size()
	return Grid{width: width, height: height}
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

func calculatePadding(dimension int, gridSize int) int {
	return (dimension - gridSize) / 2
}
func getPadding(terminal Grid) Pos {
	x := calculatePadding(terminal.width, Board.width)
	y := calculatePadding(terminal.height, Board.height)
	return Pos{x: x, y: y}
}

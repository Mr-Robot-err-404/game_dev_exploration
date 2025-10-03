package main

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
	padding  Padding
}
type Padding struct {
	top    int
	bottom int
	right  int
	left   int
}

func (gm *GameState) drawCell(pos Pos, char rune) {
	termbox.SetCell(gm.xPos(pos.x), gm.yPos(pos.y), char, termbox.ColorGreen, termbox.ColorBlack)
}
func (gm *GameState) xPos(x int) int {
	return gm.padding.left + x
}
func (gm *GameState) yPos(y int) int {
	return gm.padding.top + y
}

var Board = Grid{width: 60, height: 18}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	terminal := setup()
	padding := getPadding(terminal)
	game := GameState{board: Board, terminal: terminal, padding: padding}

	game.createGrid()

	termbox.Flush()
	time.Sleep(time.Second * 5)
}

func setup() Grid {
	width, height := termbox.Size()
	return Grid{width: width, height: height}
}

func (gm *GameState) createGrid() {
	padding := Pos{x: gm.padding.left, y: gm.padding.top}
	for i := range Board.width {
		x := padding.x + i
		y := padding.y
		termbox.SetCell(x, y, '█', termbox.ColorGreen, termbox.ColorBlack)

		y = padding.y + gm.board.height - 1
		termbox.SetCell(x, y, '█', termbox.ColorGreen, termbox.ColorBlack)
	}
	for i := range Board.height {
		y := i + padding.y
		x := 0 + padding.x
		termbox.SetCell(x, y, '█', termbox.ColorGreen, termbox.ColorBlack)

		x = pad(gm.terminal.width, padding.x)
		termbox.SetCell(x, y, '█', termbox.ColorGreen, termbox.ColorBlack)
	}
}

func calculatePadding(dimension int, gridSize int) int {
	return (dimension - gridSize) / 2
}
func getPadding(terminal Grid) Padding {
	x := calculatePadding(terminal.width, Board.width)
	y := calculatePadding(terminal.height, Board.height)
	return Padding{
		left:   x,
		top:    y,
		bottom: normalize(y),
		right:  normalize(x),
	}
}
func isEven(n int) bool {
	return n%2 == 0
}
func normalize(padding int) int {
	if isEven(padding) {
		return padding
	}
	return padding - 1
}
func pad(dimension int, padding int) int {
	pos := dimension - padding
	if isEven(dimension) {
		return pos
	}
	return pos - 1
}

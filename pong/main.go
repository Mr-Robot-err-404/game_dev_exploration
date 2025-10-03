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
type Draw struct {
	pos  Pos
	char rune
}
type Target struct {
	direction string
	value     int
}
type Coords struct {
	x Target
	y Target
}

func (gm *GameState) drawCell(coords Coords, char rune) {
	termbox.SetCell(gm.getPos(coords.x), gm.getPos(coords.y), char, termbox.ColorGreen, termbox.ColorBlack)
}
func (gm *GameState) getPos(target Target) int {
	switch target.direction {
	case "left":
		return gm.leftPos(target.value)
	case "right":
		return gm.rightPos(target.value)
	case "top":
		return gm.topPos(target.value)
	case "bottom":
		return gm.bottomPos(target.value)
	}
	return 0
}
func (gm *GameState) leftPos(x int) int {
	return gm.padding.left + x
}
func (gm *GameState) topPos(y int) int {
	return gm.padding.top + y
}
func (gm *GameState) rightPos(x int) int {
	return gm.padding.right + x
}
func (gm *GameState) bottomPos(y int) int {
	return gm.padding.bottom + y
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
	coords := Coords{
		x: Target{direction: "left"},
		y: Target{direction: "top"},
	}
	char := '█'
	for i := 0; i <= Board.width; i++ {
		coords.x.value = i
		coords.y.value = 0
		gm.drawCell(coords, char)

		coords.y.value = gm.board.height
		gm.drawCell(coords, char)
	}
	for i := 0; i <= Board.height; i++ {
		coords.x.value = 0
		coords.y.value = i
		gm.drawCell(coords, char)

		coords.x.value = gm.board.width
		gm.drawCell(coords, char)
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

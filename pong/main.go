package main

// TODO:
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

const FPS_30 = 33 * time.Millisecond

var ch = make(chan keyboardEvent)
var done = make(chan bool)

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
	player   Player
	opponent Player
}
type Player struct {
	position Pos
	movement int
}
type Draw struct {
	pos  Pos
	char rune
}

func (gm *GameState) drawPlayer(pos Pos) {
	y := pos.y / 2
	mapped := Pos{x: pos.x, y: y}

	if isEven(pos.y) {
		gm.drawCell(mapped, '█')
		gm.drawCell(Pos{x: mapped.x, y: y + 1}, '█')
		return
	}
	gm.drawCell(mapped, '▄')
	gm.drawCell(Pos{x: mapped.x, y: y + 1}, '█')
	gm.drawCell(Pos{x: mapped.x, y: y + 2}, '▀')
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
func (gm *GameState) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	gm.createBoard()
	gm.drawPlayer(gm.player.position)
	gm.drawPlayer(gm.opponent.position)
	termbox.Flush()
}
func (gm *GameState) move() {
	max := (gm.board.height - 2) * 2
	switch gm.player.movement {
	case UP:
		gm.player.position.y = dec(gm.player.position.y, 2)
	case DOWN:
		gm.player.position.y = inc(gm.player.position.y, max)
	}
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
		player: Player{
			position: Pos{x: 2, y: Board.height / 3},
			movement: STOP,
		},
		opponent: Player{
			position: Pos{x: Board.width - 2, y: Board.height / 2},
			movement: STOP,
		},
	}
	go receiveKeyboardInput(ch)
	go updateState(&game, ch, done)

	for {
		select {
		case <-done:
			return
		default:
			game.move()
			game.render()
			time.Sleep(FPS_30)
		}
	}
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
func isEven(n int) bool {
	return n%2 == 0
}

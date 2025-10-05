package main

// TODO: handle flipped logic
// ai logic

import (
	"time"

	"github.com/nsf/termbox-go"
)

const FPS_30 = 33 * time.Millisecond

const (
	TOP_LEFT int = iota
	TOP_RIGHT
	BOTTOM_LEFT
	BOTTOM_RIGHT
)
const (
	STD int = iota
	ALT
)
const (
	PLAYER_ONE int = 1 + iota
	PLAYER_TWO
)
const (
	SMALL  = 2
	NORMAL = 4
	DOUBLE = 8
)

var ch = make(chan int)
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
	board       Grid
	terminal    Grid
	padding     Pos
	player      Player
	opponent    Player
	ball        Ball
	frames      int
	active      bool
	paused      bool
	orientation int
}
type Player struct {
	position Pos
	movement int
	id       int
	size     int
}
type Ball struct {
	position Pos
	movement Movement
	maxPos   Pos
}
type Movement struct {
	north bool
	east  bool
	south bool
	west  bool
}

var Board = Grid{width: 60, height: 20}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	ort := ALT
	terminal := setup()
	padding := getPadding(terminal)
	player_pos, opponent_pos := startingPositions(ort)

	game := GameState{
		board:       Board,
		terminal:    terminal,
		padding:     padding,
		orientation: ort,
		player: Player{
			position: player_pos,
			movement: STOP,
			id:       PLAYER_ONE,
			size:     NORMAL,
		},
		opponent: Player{
			position: opponent_pos,
			movement: STOP,
			id:       PLAYER_TWO,
			size:     NORMAL,
		},
		ball: Ball{
			position: Pos{x: Board.width, y: Board.height},
			movement: Movement{west: true, north: true},
			maxPos: Pos{
				x: Board.width*2 - 1,
				y: Board.height*2 - 1,
			},
		},
	}
	go receiveKeyboardInput(ch, &game)
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
func (gm *GameState) pause() {
	gm.paused = true
}
func (gm *GameState) play() {
	gm.paused = false
}

func startingPositions(orientation int) (Pos, Pos) {
	if orientation == ALT {
		x := Board.width - 1
		return Pos{x: x, y: 4}, Pos{x: x, y: Board.height*2 - 4}
	}
	return Pos{x: 5, y: Board.height - 1}, Pos{x: Board.width*2 - 5, y: Board.height - 1}
}

func setup() Grid {
	width, height := termbox.Size()
	return Grid{width: width, height: height}
}

func playerBody(start int, size int) []int {
	body := []int{}
	for n := range size {
		body = append(body, start+n)
	}
	return body
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

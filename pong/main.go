package main

// TODO: ai logic
// scoreboard & controls text
// pause menu

import (
	"time"

	"github.com/nsf/termbox-go"
)

const FPS_30 = 33 * time.Millisecond
const FPS_60 = 17 * time.Millisecond
const FPS_90 = 12 * time.Millisecond
const FPS_120 = 8 * time.Millisecond

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

var signal = make(chan bool)

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
	player      *Player
	opponent    *Player
	ai          Ai
	ball        Ball
	frames      int
	active      bool
	paused      bool
	orientation int
	log         *Logger
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

	player, opponent := startingPlayers(ort)
	player.movement, opponent.movement = STOP, STOP
	player.id, opponent.id = PLAYER_ONE, PLAYER_TWO

	log := &Logger{ch: make(chan string, 100)}

	game := GameState{
		board:       Board,
		terminal:    terminal,
		padding:     padding,
		orientation: ort,
		player:      &player,
		opponent:    &opponent,
		ball: Ball{
			position: Pos{x: Board.width, y: Board.height},
			movement: Movement{west: true, north: true},
			maxPos: Pos{
				x: Board.width*2 - 1,
				y: Board.height*2 - 1,
			},
		},
		log: log,
		ai:  Ai{player: &opponent},
	}
	ch := make(chan int)
	mv := make(chan Mv)
	done := make(chan bool)

	go ai(&game, mv)
	go receiveKeyboardInput(ch, &game, mv)
	go updateState(&game, ch, done, mv)

	go log.init()
	log.br()
	log.msg("game started")

	ping()

	for {
		select {
		case <-done:
			return
		default:
			game.sync()
			game.move()
			game.render()
			time.Sleep(FPS_90)
		}
	}
}
func ping() {
	signal <- true
}
func (gm *GameState) sync() {
	body := playerBody(gm.ai.player.position.x, gm.ai.player.size)
	coords := gm.ai.current.target.coords

	if !gm.ai.current.has_reached && inTargetArea(coords.x, body) {
		gm.ai.current.has_reached = true
		ping()
	}
}

func (gm *GameState) pause() {
	gm.paused = !gm.paused
}
func (gm *GameState) play() {
	gm.paused = false
}

func startingPlayers(orientation int) (Player, Player) {
	if orientation == ALT {
		x := Board.width - 1
		return Player{
				position: Pos{x: x, y: 3},
				size:     DOUBLE,
			}, Player{
				position: Pos{x: x, y: Board.height*2 - 2},
				size:     DOUBLE,
			}
	}
	return Player{
			position: Pos{x: 5, y: Board.height - 1},
			size:     NORMAL,
		}, Player{
			position: Pos{x: Board.width*2 - 5, y: Board.height - 1},
			size:     NORMAL,
		}
}

func setup() Grid {
	width, height := termbox.Size()
	return Grid{width: width, height: height}
}

func playerBody(start int, size int) []int {
	body := []int{}
	end := start + size
	for n := start; n < end; n++ {
		body = append(body, n)
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

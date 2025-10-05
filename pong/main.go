package main

// TODO: ai logic

import (
	"slices"
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
	PLAYER_ONE int = 1 + iota
	PLAYER_TWO
)
const (
	SMALL  = 2
	NORMAL = 4
	LARGE  = 6
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
	board    Grid
	terminal Grid
	padding  Pos
	player   Player
	opponent Player
	ball     Ball
	frames   int
	active   bool
	paused   bool
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

func (gm *GameState) drawPlayer(pos Pos, id int) {
	y := pos.y / 2
	x := pos.x / 2
	mapped := Pos{x: x, y: y}

	if isEven(pos.y) {
		if id == PLAYER_ONE {
			gm.drawCell(mapped, '▐')
			gm.drawCell(Pos{x: mapped.x, y: y + 1}, '▐')
			return
		}
		gm.drawCell(mapped, '▌')
		gm.drawCell(Pos{x: mapped.x, y: y + 1}, '▌')
		return
	}
	if id == PLAYER_ONE {
		gm.drawCell(mapped, '▗')
		gm.drawCell(Pos{x: mapped.x, y: y + 1}, '▐')
		gm.drawCell(Pos{x: mapped.x, y: y + 2}, '▝')
		return
	}
	gm.drawCell(mapped, '▖')
	gm.drawCell(Pos{x: mapped.x, y: y + 1}, '▌')
	gm.drawCell(Pos{x: mapped.x, y: y + 2}, '▘')
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
func (gm *GameState) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	gm.createBoard()
	gm.drawPlayer(gm.player.position, PLAYER_ONE)
	gm.drawPlayer(gm.opponent.position, PLAYER_TWO)
	gm.drawBall()

	termbox.Flush()
}
func (gm *GameState) move() {
	if gm.paused {
		return
	}
	gm.collisions()
	gm.movePlayer()
	gm.moveBall()
}

func (gm *GameState) movePlayer() {
	min := 2
	max := (gm.board.height - 2) * 2
	switch gm.player.movement {
	case UP:
		y := dec(gm.player.position.y, 2)
		gm.player.position.y = y
		if y == min {
			gm.player.movement = STOP
		}
	case DOWN:
		y := inc(gm.player.position.y, max)
		gm.player.position.y = y
		if y == max {
			gm.player.movement = STOP
		}
	}
}

func (gm *GameState) collisions() {
	gm.wallCollision()
	gm.playerCollision(gm.player)
	gm.playerCollision(gm.opponent)
}
func (gm *GameState) wallCollision() {
	min, max := 2, gm.ball.maxPos
	mv, pos := gm.ball.movement, gm.ball.position

	if pos.x == max.x && mv.east {
		gm.invertXMovement()
	}
	if pos.x == min && mv.west {
		gm.invertXMovement()
	}
	if pos.y == max.y && mv.south {
		gm.invertYMovement()
	}
	if pos.y == min && mv.north {
		gm.invertYMovement()
	}
}

func (gm *GameState) playerCollision(player Player) {
	ball := gm.ball.position
	pos := player.position

	x := pos.x + 1
	if player.id == PLAYER_TWO {
		x = pos.x - 1
	}
	if ball.x != x {
		return
	}
	body := playerBody(pos.y, player.size)

	if !slices.Contains(body, ball.y) {
		return
	}
	defer gm.invertXMovement()

	switch player.movement {
	case UP:
		gm.moveBallNorth()
	case DOWN:
		gm.moveBallSouth()
	}
}

func (gm *GameState) moveBall() {
	mv := gm.ball.movement
	maxH := gm.board.height*2 - 1
	maxW := gm.board.width*2 - 1

	if mv.south {
		gm.ball.position.y = inc(gm.ball.position.y, maxH)
	}
	if mv.north {
		gm.ball.position.y = dec(gm.ball.position.y, 2)
	}
	if mv.east {
		gm.ball.position.x = inc(gm.ball.position.x, maxW)
	}
	if mv.west {
		gm.ball.position.x = dec(gm.ball.position.x, 2)
	}
}

func (gm *GameState) moveBallNorth() {
	gm.ball.movement.south = false
	gm.ball.movement.north = true
}
func (gm *GameState) moveBallSouth() {
	gm.ball.movement.north = false
	gm.ball.movement.south = true
}
func (gm *GameState) invertXMovement() {
	if gm.ball.movement.east {
		gm.ball.movement.east = false
		gm.ball.movement.west = true
		return
	}
	gm.ball.movement.west = false
	gm.ball.movement.east = true
}
func (gm *GameState) invertYMovement() {
	if gm.ball.movement.north {
		gm.ball.movement.north = false
		gm.ball.movement.south = true
		return
	}
	gm.ball.movement.south = false
	gm.ball.movement.north = true
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

var Board = Grid{width: 60, height: 20}

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
			position: Pos{x: 5, y: Board.height - 1},
			movement: STOP,
			id:       PLAYER_ONE,
			size:     NORMAL,
		},
		opponent: Player{
			position: Pos{x: Board.width*2 - 5, y: Board.height - 1},
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
	go receiveKeyboardInput(ch)
	go updateState(&game, ch, done)

	for {
		select {
		case <-done:
			return
		default:
			game.render()
			game.move()
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
func playerBody(y int, size int) []int {
	body := []int{}
	for n := range size {
		body = append(body, y+n)
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

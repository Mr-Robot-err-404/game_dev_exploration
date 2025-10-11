package main

import (
	"fmt"
	"slices"
)

// TODO:
// handle side collisions

const (
	NORTH int = iota
	EAST
	SOUTH
	WEST
)
const (
	START int = iota
	BOUNCE
	TARGET_AREA
)

var ctx_map = map[int]string{
	START:       "start",
	BOUNCE:      "bounce",
	TARGET_AREA: "target_area",
}

var input_map = map[int]int{
	NORTH: UP,
	EAST:  RIGHT,
	SOUTH: DOWN,
	WEST:  LEFT,
}

type Target struct {
	coords Pos
	active bool
}
type CurrentTarget struct {
	target      Target
	has_reached bool
}
type Ai struct {
	player    *Player
	home      Target
	intercept Target
	current   CurrentTarget
	input     chan<- Mv
	log       *Logger
}

func ai(gm *GameState) {
	gm.ai.home = Target{
		coords: Pos{
			y: gm.ai.player.position.y,
			x: 5},
		active: true,
	}
	gm.ai.current = CurrentTarget{has_reached: false, target: gm.ai.home}
	for {
		ctx := <-signal

		ai := gm.ai
		player := ai.player

		position := player.position
		body := playerBody(position.x, player.size)

		gm.log.br()
		gm.log.msg(ctx_map[ctx])
		gm.log.msg(fmt.Sprintf("%d:%d", body[0], body[len(body)-1]))

		if ctx == BOUNCE || ctx == START {
			ai.bounce(gm.log, gm.ball, gm.board)
			continue
		}
		if ctx == TARGET_AREA {
			ai.targetArea(gm.log)
		}
	}
}

func (ai *Ai) targetArea(log *Logger) {
	if ai.home.active {
		ai.atHome(log)
		return
	}
	ai.atIntercept(log)
}

func (ai *Ai) bounce(log *Logger, ball Ball, board Grid) {
	if ball.movement.south {
		end := ai.player.position.y - 1
		pos, _ := walk(ball, end, 0, board)
		pos.y++

		ai.assignIntercept(pos)
		ai.targetIntercept(log)
		return
	}
	ai.assignHome(ai.home.coords)
	ai.targetHome(log)
}

func (ai *Ai) assignHome(pos Pos) {
	ai.home.coords = pos
	ai.home.active = true
	ai.current.target = ai.home
	ai.current.has_reached = false
	ai.intercept.active = false
}

func (ai *Ai) assignIntercept(pos Pos) {
	ai.intercept.coords = pos
	ai.intercept.active = true
	ai.current.target = ai.intercept
	ai.current.has_reached = false
	ai.home.active = false
}

func (ai *Ai) targetIntercept(log *Logger) {
	if ai.current.has_reached {
		ai.atIntercept(log)
		return
	}
	ai.goToIntercept(log)
}

func (ai *Ai) goToIntercept(log *Logger) {
	compass := findHome(ai.player.position, ai.intercept.coords)
	move := input_map[compass]

	log.msg(fmt.Sprintf("PING: current pos x:%d, intercept x:%d", ai.player.position.x, ai.intercept.coords.x))

	if move != ai.player.movement {
		ai.input <- Mv{event: move, player_id: ai.player.id}
		log.msg(fmt.Sprintf("PING: sent move command %d", move))
		return
	}
	log.msg("PING: already moving towards intercept")
}

func (ai *Ai) atIntercept(log *Logger) {
	log.msg("PING: at intercept, stopping")
	ai.intercept.active = false

	if ai.player.movement != STOP {
		ai.input <- Mv{event: STOP, player_id: ai.player.id}
		log.msg("PING: sent STOP command")
	}
}

func (ai *Ai) targetHome(log *Logger) {
	if ai.current.has_reached {
		ai.atHome(log)
		return
	}
	ai.goHome(log)
}

func (ai *Ai) goHome(log *Logger) {
	compass := findHome(ai.player.position, ai.home.coords)
	move := input_map[compass]

	log.msg(fmt.Sprintf("PING: x:%d, home:%d", ai.player.position.x, ai.home.coords.x))

	if move != ai.player.movement {
		ai.input <- Mv{event: move, player_id: ai.player.id}
		log.msg(fmt.Sprintf("PING: sent move command %d", move))
	} else {
		log.msg("PING: already going home")
	}
}

func (ai *Ai) atHome(log *Logger) {
	log.msg("PING: at home, stopping")
	ai.home.active = false

	if ai.player.movement != STOP {
		ai.input <- Mv{event: STOP, player_id: ai.player.id}
		log.msg("PING: sent STOP command")
	}
}

func walk(ball Ball, end int, steps int, board Grid) (Pos, int) {
	if ball.position.y == end {
		return ball.position, steps
	}
	return walk(stepForward(ball, board), end, steps+1, board)
}

func inTargetArea(axis int, body []int) bool {
	return slices.Contains(body, axis)
}
func findHome(pos Pos, home Pos) int {
	if pos.x < home.x {
		return EAST
	}
	if pos.x > home.x {
		return WEST
	}
	if pos.y < home.y {
		return NORTH
	}
	return SOUTH
}

func stepForward(ball Ball, board Grid) Ball {
	mv := ball.movement
	maxH := board.height*2 - 1
	maxW := board.width*2 - 1
	min := 2

	collision(&ball, maxW-1, min+1)

	if mv.south {
		ball.position.y = inc(ball.position.y, maxH)
	}
	if mv.north {
		ball.position.y = dec(ball.position.y, min)
	}
	if mv.east {
		ball.position.x = inc(ball.position.x, maxW)
	}
	if mv.west {
		ball.position.x = dec(ball.position.x, min)
	}
	return ball
}

func collision(ball *Ball, max int, min int) {
	if ball.movement.east && ball.position.x == max {
		ball.movement.east = false
		ball.movement.west = true
		return
	}
	if ball.movement.west && ball.position.x == min {
		ball.movement.west = false
		ball.movement.east = true
	}
}

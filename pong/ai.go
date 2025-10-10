package main

import (
	"fmt"
	"slices"
)

// TODO: ping:
// change of direction
// if in range for target

const (
	NORTH int = iota
	EAST
	SOUTH
	WEST
)

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
}

func ai(gm *GameState, input chan<- Mv) {
	gm.ai.home = Target{
		coords: Pos{
			y: gm.ai.player.position.y,
			x: 15},
		active: true,
	}
	gm.ai.current = CurrentTarget{has_reached: false, target: gm.ai.home}
	for {
		<-signal
		ai := gm.ai
		player, intercept := ai.player, ai.intercept

		position := player.position
		body := playerBody(position.x, player.size)

		gm.log.br()
		gm.log.msg(fmt.Sprintf("body -> %d:%d", body[0], body[len(body)-1]))

		if ai.home.active {
			ai.targetHome(input, gm.log)
			continue
		}
		continue
		end := player.position.y - 1
		pos, _ := walk(gm.ball, end, 0, gm.board)

		if !intercept.active {
			pos.y++
			gm.ai.intercept.coords = pos
		}
	}
}

func (ai *Ai) targetHome(input chan<- Mv, log *Logger) {
	if ai.current.has_reached {
		ai.atHome(input, log)
		return
	}
	ai.goHome(input, log)
}

func (ai *Ai) goHome(input chan<- Mv, log *Logger) {
	compass := findHome(ai.player.position, ai.home.coords)
	move := input_map[compass]

	log.msg(fmt.Sprintf("PING: need to move compass:%d move:%d", compass, move))
	log.msg(fmt.Sprintf("PING: current pos x:%d, target x:%d", ai.player.position.x, ai.home.coords.x))

	if move != ai.player.movement {
		input <- Mv{event: move, player_id: ai.player.id}
		log.msg(fmt.Sprintf("PING: sent move command %d", move))
	} else {
		log.msg("PING: already moving in correct direction")
	}
}

func (ai *Ai) atHome(input chan<- Mv, log *Logger) {
	log.msg("PING: at home, stopping")
	ai.home.active = false

	if ai.player.movement != STOP {
		input <- Mv{event: STOP, player_id: ai.player.id}
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

	if mv.south {
		ball.position.y = inc(ball.position.y, maxH)
	}
	if mv.north {
		ball.position.y = dec(ball.position.y, 2)
	}
	if mv.east {
		ball.position.x = inc(ball.position.x, maxW)
	}
	if mv.west {
		ball.position.x = dec(ball.position.x, 2)
	}
	return ball
}

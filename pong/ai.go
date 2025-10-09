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
type Strategy struct {
	prev   Movement
	target Target
}
type Ai struct {
	strat  Strategy
	player *Player
	home   Target
}

func ai(gm *GameState, input chan<- Mv) {
	gm.ai.home = Target{coords: Pos{
		y: gm.ai.player.position.y,
		x: gm.board.width,
	}}
	for {
		<-signal
		player, strat, home := gm.ai.player, gm.ai.strat, gm.ai.home
		position, mv := player.position, player.movement
		body := playerBody(position.x, player.size)

		gm.log.msg(fmt.Sprintf("size -> %d:%d", position.x, player.size))
		gm.log.msg(fmt.Sprintf("body -> %d:%d", body[0], body[len(body)-1]))

		if gm.ball.movement.north {
			if isHome(position.x, body) && mv != STOP {
				input <- Mv{event: STOP, player_id: player.id}
				gm.ai.home.active = false
				return
			}
			gm.ai.home.active = true

			compass := findHome(player.position, home.coords)
			move := input_map[compass]

			gm.log.msg(fmt.Sprintf("compass: %d", compass))
			gm.log.msg(fmt.Sprintf("move: %d", move))
			gm.log.msg(fmt.Sprintf("player_mv: %d", player.movement))

			if move != player.movement {
				input <- Mv{event: move, player_id: player.id}
			}
			return
		}
		end := player.position.y - 1
		pos, steps := walk(gm.ball, end, 0, gm.board)

		if !strat.target.active {
			pos.y++
			strat.target.coords = pos
			strat.target.active = true
			gm.opponent.position = pos

			gm.log.msg(fmt.Sprintf("target -> %d:%d", pos.x, pos.y))
			gm.log.msg(fmt.Sprintf("steps -> %d", steps))
		}
	}
}

func walk(ball Ball, end int, steps int, board Grid) (Pos, int) {
	if ball.position.y == end {
		return ball.position, steps
	}
	return walk(stepForward(ball, board), end, steps+1, board)
}

func isHome(axis int, body []int) bool {
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

package main

import "fmt"

type Target struct {
	coords Pos
	active bool
}
type Strategy struct {
	prev   Movement
	target Target
	home   Pos
}

func ai(gm *GameState, ping <-chan bool) {
	target := Target{active: false}
	for {
		<-ping
		end := gm.opponent.position.y - 1
		pos, steps := walk(gm.ball, end, 0, gm.board)

		if !target.active {
			pos.y++
			target.coords = pos
			target.active = true
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

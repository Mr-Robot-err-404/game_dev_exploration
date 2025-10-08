package main

import (
	"github.com/nsf/termbox-go"
)

type Mv struct {
	event     int
	player_id int
}

const (
	UP int = iota
	DOWN
	LEFT
	RIGHT
	STOP
	PAUSE
	END
)

type Rcv struct {
	game *GameState
	ch   chan int
	done chan bool
	mv   chan Mv
}

func updateState(rcv Rcv) {
	for {
		select {
		case mv := <-rcv.mv:
			rcv.game.play()
			switch mv.event {
			case UP:
				assignMovement(rcv.game, mv.player_id, UP)
			case DOWN:
				assignMovement(rcv.game, mv.player_id, DOWN)
			case LEFT:
				assignMovement(rcv.game, mv.player_id, LEFT)
			case RIGHT:
				assignMovement(rcv.game, mv.player_id, RIGHT)
			case STOP:
				assignMovement(rcv.game, mv.player_id, STOP)
			}
		case n := <-rcv.ch:
			switch n {
			case PAUSE:
				rcv.game.pause()
			case END:
				done <- true
			}
		}

	}
}
func receiveKeyboardInput(ch chan<- int, gm *GameState, mv chan<- Mv) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'j':
				if gm.orientation == ALT {
					mv <- Mv{player_id: PLAYER_ONE, event: LEFT}
					continue
				}
				mv <- Mv{player_id: PLAYER_ONE, event: DOWN}
			case 'k':
				if gm.orientation == ALT {
					mv <- Mv{player_id: PLAYER_ONE, event: RIGHT}
					continue
				}
				mv <- Mv{player_id: PLAYER_ONE, event: UP}
			case 'f':
				if gm.orientation == ALT {
					mv <- Mv{player_id: PLAYER_TWO, event: RIGHT}
					continue
				}
				mv <- Mv{player_id: PLAYER_TWO, event: UP}
			case 'd':
				if gm.orientation == ALT {
					mv <- Mv{player_id: PLAYER_TWO, event: LEFT}
					continue
				}
				mv <- Mv{player_id: PLAYER_TWO, event: DOWN}

			case 'q':
				ch <- PAUSE
			}
			switch ev.Key {
			case termbox.KeyEsc:
				ch <- END
			case termbox.KeySpace:
				mv <- Mv{player_id: PLAYER_ONE, event: STOP}
			case termbox.KeyEnter:
				mv <- Mv{player_id: PLAYER_TWO, event: STOP}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
func inc(num int, max int) int {
	n := num + 1
	if n >= max {
		return max
	}
	return n
}
func dec(num int, min int) int {
	n := num - 1
	if n < min {
		return min
	}
	return n
}
func assignMovement(gm *GameState, id int, mv int) {
	if id == PLAYER_ONE {
		gm.player.movement = mv
	} else {
		gm.opponent.movement = mv
	}
}

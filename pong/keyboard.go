package main

import "github.com/nsf/termbox-go"

type keyboardEvent struct {
	event int
	key   termbox.Key
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

func updateState(game *GameState, ch chan int, done chan bool) {
	for {
		rcv := <-ch
		if rcv != PAUSE {
			game.play()
		}
		switch rcv {
		case UP:
			game.player.movement = UP
		case DOWN:
			game.player.movement = DOWN
		case LEFT:
			game.player.movement = LEFT
		case RIGHT:
			game.player.movement = RIGHT
		case STOP:
			game.player.movement = STOP
		case PAUSE:
			game.pause()
		case END:
			done <- true
		}
	}
}
func receiveKeyboardInput(ch chan<- int, gm *GameState) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'j':
				if gm.orientation == ALT {
					ch <- LEFT
					continue
				}
				ch <- DOWN
			case 'k':
				if gm.orientation == ALT {
					ch <- RIGHT
					continue
				}
				ch <- UP
			case 'q':
				ch <- PAUSE
			}
			switch ev.Key {
			case termbox.KeyEsc:
				ch <- END
			case termbox.KeySpace:
				ch <- STOP
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

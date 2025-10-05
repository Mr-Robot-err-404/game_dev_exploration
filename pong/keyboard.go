package main

import "github.com/nsf/termbox-go"

type keyboardEvent struct {
	event int
	key   termbox.Key
}

const (
	UP int = iota
	DOWN
	STOP
	PAUSE
	END
)

func updateState(game *GameState, ch chan int, done chan bool) {
	for {
		rcv := <-ch
		switch rcv {
		case UP:
			game.player.movement = UP
			game.paused = false
		case DOWN:
			game.player.movement = DOWN
			game.paused = false
		case STOP:
			game.player.movement = STOP
			game.paused = false
		case PAUSE:
			game.paused = !game.paused
		case END:
			done <- true
		}
	}
}
func receiveKeyboardInput(ch chan<- int) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'j':
				ch <- DOWN
			case 'k':
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

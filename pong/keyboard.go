package main

import "github.com/nsf/termbox-go"

type keyboardEvent struct {
	event int
	key   termbox.Key
}

const (
	UP int = 1 + iota
	DOWN
	STOP
	END
)

func updateState(game *GameState, ch chan keyboardEvent, done chan bool) {
	for {
		rcv := <-ch
		switch rcv.event {
		case UP:
			game.player.movement = UP
		case DOWN:
			game.player.movement = DOWN
		case STOP:
			game.player.movement = STOP
		case END:
			done <- true
		}
	}
}
func receiveKeyboardInput(ch chan<- keyboardEvent) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'j':
				ch <- keyboardEvent{event: DOWN, key: termbox.Key(ev.Ch)}
			case 'k':
				ch <- keyboardEvent{event: UP, key: termbox.Key(ev.Ch)}
			}
			switch ev.Key {
			case termbox.KeyEsc:
				ch <- keyboardEvent{event: END, key: ev.Key}
			case termbox.KeySpace:
				ch <- keyboardEvent{event: STOP, key: ev.Key}
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

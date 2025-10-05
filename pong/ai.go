package main

import "time"

func ai(gm *GameState, ch chan<- int) {
	ticker := time.NewTicker(100 * time.Millisecond)

	defer ticker.Stop()
	for {
		<-ticker.C
		if !gm.active {
			continue
		}
	}
}

func walk() {}

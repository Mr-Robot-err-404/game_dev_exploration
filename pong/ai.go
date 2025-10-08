package main

// TODO: project direction of ball
// find end coordinate

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
	for {
		<-ping
	}
}

func projectPath() {}

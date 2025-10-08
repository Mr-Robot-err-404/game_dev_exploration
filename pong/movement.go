package main

import (
	"slices"
)

func (gm *GameState) move() {
	if gm.paused {
		return
	}
	gm.movePlayer(&gm.player)
	gm.movePlayer(&gm.opponent)

	if isEven(gm.frames) {
		gm.collisions()
		gm.moveBall()
	}
	gm.frames++
}

func (gm *GameState) movePlayer(player *Player) {
	min := 2
	maxH := (gm.board.height - 2) * 2
	maxW := (gm.board.width - 2) * 2

	switch player.movement {
	case UP:
		y := dec(player.position.y, 2)
		player.position.y = y
		if y == min {
			player.movement = STOP
		}
	case DOWN:
		y := inc(player.position.y, maxH)
		player.position.y = y
		if y == maxH {
			player.movement = STOP
		}
	case LEFT:
		x := dec(player.position.x, min)
		player.position.x = x
		if x == min {
			player.movement = STOP
		}
	case RIGHT:
		x := inc(player.position.x, maxW-4)
		player.position.x = x
		if x == maxW {
			player.movement = STOP
		}
	}
}

func (gm *GameState) collisions() {
	gm.wallCollision()

	if gm.orientation == ALT {
		gm.altPlayerCollision(gm.player)
		gm.altPlayerCollision(gm.opponent)
		return
	}
	gm.playerCollision(gm.player)
	gm.playerCollision(gm.opponent)
}
func (gm *GameState) wallCollision() {
	min, max := 2, gm.ball.maxPos
	mv, pos := gm.ball.movement, gm.ball.position

	if pos.x == max.x && mv.east {
		gm.invertXMovement()
	}
	if pos.x == min && mv.west {
		gm.invertXMovement()
	}
	if pos.y == max.y && mv.south {
		gm.invertYMovement()
	}
	if pos.y == min && mv.north {
		gm.invertYMovement()
	}
}

func (gm *GameState) altPlayerCollision(player Player) {
	ball := gm.ball.position
	pos := player.position

	y := pos.y + 1
	if player.id == PLAYER_TWO {
		y = pos.y - 1
	}
	if ball.y != y {
		return
	}
	body := playerBody(pos.x, player.size)

	if !slices.Contains(body, ball.x) {
		return
	}
	defer gm.invertYMovement()

	switch player.movement {
	case LEFT:
		gm.moveBallWest()
	case RIGHT:
		gm.moveBallEast()
	}
}

func (gm *GameState) playerCollision(player Player) {
	ball := gm.ball.position
	pos := player.position

	x := pos.x + 1
	if player.id == PLAYER_TWO {
		x = pos.x - 1
	}
	if ball.x != x {
		return
	}
	body := playerBody(pos.y, player.size)

	if !slices.Contains(body, ball.y) {
		return
	}
	defer gm.invertXMovement()

	switch player.movement {
	case UP:
		gm.moveBallNorth()
	case DOWN:
		gm.moveBallSouth()
	}
}

func (gm *GameState) moveBall() {
	mv := gm.ball.movement
	maxH := gm.board.height*2 - 1
	maxW := gm.board.width*2 - 1

	if mv.south {
		gm.ball.position.y = inc(gm.ball.position.y, maxH)
	}
	if mv.north {
		gm.ball.position.y = dec(gm.ball.position.y, 2)
	}
	if mv.east {
		gm.ball.position.x = inc(gm.ball.position.x, maxW)
	}
	if mv.west {
		gm.ball.position.x = dec(gm.ball.position.x, 2)
	}
}

func (gm *GameState) moveBallNorth() {
	gm.ball.movement.south = false
	gm.ball.movement.north = true
}
func (gm *GameState) moveBallSouth() {
	gm.ball.movement.north = false
	gm.ball.movement.south = true
}
func (gm *GameState) moveBallWest() {
	gm.ball.movement.east = false
	gm.ball.movement.west = true
}
func (gm *GameState) moveBallEast() {
	gm.ball.movement.west = false
	gm.ball.movement.east = true
}
func (gm *GameState) invertXMovement() {
	if gm.ball.movement.east {
		gm.ball.movement.east = false
		gm.ball.movement.west = true
		return
	}
	gm.ball.movement.west = false
	gm.ball.movement.east = true
}
func (gm *GameState) invertYMovement() {
	if gm.ball.movement.north {
		gm.ball.movement.north = false
		gm.ball.movement.south = true
		return
	}
	gm.ball.movement.south = false
	gm.ball.movement.north = true
}

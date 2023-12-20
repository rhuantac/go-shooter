package main

import "github.com/ByteArena/box2d"


type Command string

const (
	MoveLeft Command = "LEFT"
	MoveUp Command = "UP"
	MoveRight Command = "RIGHT"
	MoveDown Command = "DOWN"
)

const MoveSpeed = 10.0
type Action struct {
	command Command
	player string
}



func ProcessWorld(world *box2d.B2World, playerStore map[string]*box2d.B2Body, action Action) error{
	player := playerStore[action.player]
	if action.command == MoveUp {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: MoveSpeed}, true)
	} else if action.command == MoveLeft {
		player.ApplyForceToCenter(box2d.B2Vec2{X: MoveSpeed * -1, Y: 0.0}, true)
	} else if action.command == MoveDown {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: MoveSpeed * -1}, true)
	} else if action.command == MoveRight{
		player.ApplyForceToCenter(box2d.B2Vec2{X: MoveSpeed, Y: 0.0}, true)
	}
	
	for i := 0; i < 60; i++ {
		world.Step(1.0 / 60.0, 6, 2)
	}
	return nil
}
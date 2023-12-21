package main

import "github.com/ByteArena/box2d"


type InputType string

const (
	MoveLeft InputType = "LEFT"
	MoveUp InputType = "UP"
	MoveRight InputType = "RIGHT"
	MoveDown InputType = "DOWN"
)

const MoveSpeed = 10.0
type Action struct {
	input InputType
	player string
}



func ProcessWorld(world *box2d.B2World, playerStore map[string]*box2d.B2Body, action Action) error{
	player := playerStore[action.player]
	if action.input == MoveUp {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: MoveSpeed}, true)
	} else if action.input == MoveLeft {
		player.ApplyForceToCenter(box2d.B2Vec2{X: MoveSpeed * -1, Y: 0.0}, true)
	} else if action.input == MoveDown {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: MoveSpeed * -1}, true)
	} else if action.input == MoveRight{
		player.ApplyForceToCenter(box2d.B2Vec2{X: MoveSpeed, Y: 0.0}, true)
	}
	
	for i := 0; i < 60; i++ {
		world.Step(1.0 / 60.0, 6, 2)
	}
	return nil
}
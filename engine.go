package main

import "github.com/ByteArena/box2d"

type Action struct {
	action string
	player string
}

func ProcessWorld(world *box2d.B2World, playerStore map[string]*box2d.B2Body, action Action) error{
	player := playerStore[action.player]
	if action.action == "UP" {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: 10.0}, true)
	} else if action.action == "LEFT" {
		player.ApplyForceToCenter(box2d.B2Vec2{X: -10.0, Y: 0.0}, true)
	} else if action.action == "DOWN" {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: -10.0}, true)
	} else {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 10.0, Y: 0.0}, true)
	}
	
	for i := 0; i < 60; i++ {
		world.Step(1.0 / 60.0, 6, 2)
	}
	return nil
}
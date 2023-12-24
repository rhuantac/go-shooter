package main

import "github.com/ByteArena/box2d"


type InputType string

const (
	MoveLeft InputType = "LEFT"
	MoveUp InputType = "UP"
	MoveRight InputType = "RIGHT"
	MoveDown InputType = "DOWN"
)

const moveSpeed = 10.0
type Action struct {
	input InputType
	player string
}
type Processor interface {
	Process(playerStore map[string]*box2d.B2Body, action Action) error
	CreateCharacter() *box2d.B2Body
}
type WorldProcessor struct {
	World *box2d.B2World
}

func (e WorldProcessor) Process(playerStore map[string]*box2d.B2Body, action Action) error{
	player := playerStore[action.player]
	if action.input == MoveUp {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: moveSpeed}, true)
	} else if action.input == MoveLeft {
		player.ApplyForceToCenter(box2d.B2Vec2{X: moveSpeed * -1, Y: 0.0}, true)
	} else if action.input == MoveDown {
		player.ApplyForceToCenter(box2d.B2Vec2{X: 0.0, Y: moveSpeed * -1}, true)
	} else if action.input == MoveRight{
		player.ApplyForceToCenter(box2d.B2Vec2{X: moveSpeed, Y: 0.0}, true)
	}
	
	for i := 0; i < 60; i++ {
		e.World.Step(1.0 / 60.0, 6, 2)
	}
	return nil
}

func (e WorldProcessor) CreateCharacter() *box2d.B2Body{
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(2.0, 2.0)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = false	
	bd.AllowSleep = false
	
	character := e.World.CreateBody(&bd)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(0.20, 0.20)
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	character.CreateFixtureFromDef(&fd)
	return character
}

func setupWorld() box2d.B2World{
	gravity := box2d.MakeB2Vec2(0, 0)
	return box2d.MakeB2World(gravity)
}

func NewProcessor() WorldProcessor {
	world := setupWorld()
	return WorldProcessor{World: &world}
}
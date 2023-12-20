package main

import (
	"math"
	"testing"

	"github.com/ByteArena/box2d"
)

func almostEqual(a, b float64) bool{
	var tolerance float64 = 0.1
	if a == b {
		return true
	}

	d := math.Abs(a - b)

	if b == 0 {
        return d < tolerance
    }
    return (d / math.Abs(b)) < tolerance
}

func setupWorld() (world box2d.B2World, character *box2d.B2Body) {
	//Setup world
	gravity := box2d.MakeB2Vec2(0, 0)
	world = box2d.MakeB2World(gravity)
	
	//Setup Character
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(2.0, 2.0)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = false	
	bd.AllowSleep = false
	character = world.CreateBody(&bd)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(0.20, 0.20)
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	character.CreateFixtureFromDef(&fd)
	return
}

func TestActions(t *testing.T) {
	t.Run("player goes up on action UP", func(t *testing.T) {
		world, character := setupWorld()
		playerName :=  "John"
		action := Action{action: "UP", player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character
		//Acions to perform
		err := ProcessWorld(&world, playerStore, action)

		if err != nil {
			t.Fatalf("Error processing world, '%v'", action)
		}

		if !almostEqual(character.GetPosition().X, 2.0) {
			t.Errorf("player '%s' moved horizontally while going up", playerName)
		}

		if !almostEqual(character.GetPosition().Y, 3.0) {
			t.Errorf("player '%s' didn't go up correctly. Got %.1f expected %.1f", playerName, character.GetPosition().Y, 3.0)
		}

	})

	t.Run("player goes left on action LEFT", func(t *testing.T) {
		world, character := setupWorld()
		playerName :=  "John"
		action := Action{action: "LEFT", player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character
		//Acions to perform
		err := ProcessWorld(&world, playerStore, action)

		if err != nil {
			t.Fatalf("Error processing world, '%v'", action)
		}

		if !almostEqual(character.GetPosition().X, 1.0) {
			t.Errorf("player '%s' didn't go left correctly. Got %.1f expected %.1f", playerName, character.GetPosition().X, 1.0)
		}

		if !almostEqual(character.GetPosition().Y, 2.0) {
			t.Errorf("player '%s' moved vertically while going left", playerName)
		}
		
	})
}
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

func TestActions(t *testing.T) {
	t.Run("player goes UP", func(t *testing.T) {
		expectedX := 2.0
		expectedY := 3.0
		worldProcessor := NewProcessor()
		character := worldProcessor.CreateCharacter()
		playerName :=  "John"
		action := Action{input: MoveUp, player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character

		worldProcessor.Process(playerStore, action)

		if !almostEqual(character.GetPosition().Y, expectedY) {
			t.Errorf("player '%s' didn't go up correctly. Got %.1f expected %.1f", playerName, character.GetPosition().Y, expectedY)
		}

		if !almostEqual(character.GetPosition().X, expectedX) {
			t.Errorf("player '%s' moved horizontally while going up", playerName)
		}
	})

	t.Run("player goes LEFT", func(t *testing.T) {
		expectedX := 1.0
		expectedY := 2.0
		worldProcessor := NewProcessor()
		character := worldProcessor.CreateCharacter()
		playerName :=  "John"
		action := Action{input: MoveLeft, player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character

		worldProcessor.Process(playerStore, action)

		if !almostEqual(character.GetPosition().X, expectedX) {
			t.Errorf("player '%s' didn't go left correctly. Got %.1f expected %.1f", playerName, character.GetPosition().X, expectedX)
		}

		if !almostEqual(character.GetPosition().Y, expectedY) {
			t.Errorf("player '%s' moved vertically while going left", playerName)
		}		
	})

	t.Run("player goes DOWN", func(t *testing.T) {
		expectedX := 2.0
		expectedY := 1.0
		worldProcessor := NewProcessor()
		character := worldProcessor.CreateCharacter()
		playerName :=  "John"
		action := Action{input: MoveDown, player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character

		worldProcessor.Process(playerStore, action)

		if !almostEqual(character.GetPosition().Y, expectedY) {
			t.Errorf("player '%s' didn't go down correctly. Got %.1f expected %.1f", playerName, character.GetPosition().Y, expectedY)
		}

		if !almostEqual(character.GetPosition().X, expectedX) {
			t.Errorf("player '%s' moved horizontally while going down", playerName)
		}		
	})

	t.Run("player goes RIGHT", func(t *testing.T) {
		expectedX := 3.0
		expectedY := 2.0
		worldProcessor := NewProcessor()
		character := worldProcessor.CreateCharacter()
		playerName :=  "John"
		action := Action{input: MoveRight, player: playerName}
		playerStore := make(map[string] *box2d.B2Body)
		playerStore[playerName] = character

		worldProcessor.Process(playerStore, action)

		if !almostEqual(character.GetPosition().X, expectedX) {
			t.Errorf("player '%s' didn't go right correctly. Got %.1f expected %.1f", playerName, character.GetPosition().X, expectedX)
		}

		if !almostEqual(character.GetPosition().Y, expectedY) {
			t.Errorf("player '%s' moved horizontally while going down", playerName)
		}		
	})
}
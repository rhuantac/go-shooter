package main

import (
	"testing"

	"github.com/ByteArena/box2d"
)

type StubWorldProcessor struct {
	createCharacterCalls int
	processCalls int
}

func (w *StubWorldProcessor) Process(playerStore map[string]*box2d.B2Body, action Action) error{
	w.processCalls += 1
	return nil	
}

func (w *StubWorldProcessor) CreateCharacter() *box2d.B2Body{
	w.createCharacterCalls += 1
	return nil
}
func TestGameManager(t *testing.T) {
	t.Run("Test that one character is created on new game", func(t *testing.T) {
		processor := StubWorldProcessor{}
		CreateNewGame(&processor)
		if processor.createCharacterCalls != 1 {
			t.Errorf("got %d calls on CreateCharacter want %d", processor.createCharacterCalls, 1)
		}
	})

	t.Run("Test that World is processed 60 ticks per call", func (t *testing.T) {
		processor := StubWorldProcessor{}
		playerStore := make(PlayerStore)
		CreateNewGame(&processor)
		gameLoop(&processor, playerStore)
		if (processor.processCalls != 60) {
			t.Errorf("got %d calls on gameLoop want %d", processor.processCalls, 60)
		}
	})
}
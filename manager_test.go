package main

import (
	"testing"
	"time"

	"github.com/ByteArena/box2d"
)

type stubWorldProcessor struct {
	createCharacterCalls int
	processCalls         int
}

func (w *stubWorldProcessor) Process(playerStore map[string]*box2d.B2Body, action []Action) error {
	w.processCalls += 1
	return nil
}

func (w *stubWorldProcessor) CreateCharacter() *box2d.B2Body {
	w.createCharacterCalls += 1
	return nil
}

func TestGameManager(t *testing.T) {
	t.Run("one character is created on new game", func(t *testing.T) {
		processor := stubWorldProcessor{}
		CreateNewGame(&processor)
		if processor.createCharacterCalls != 1 {
			t.Errorf("got %d calls on CreateCharacter want %d", processor.createCharacterCalls, 1)
		}
	})

	t.Run("world is processed 60 ticks per call", func(t *testing.T) {
		processor := stubWorldProcessor{}
		playerStore := make(PlayerStore)
		queueProcess(&processor, playerStore, &ActionQueue{})
		if processor.processCalls != 60 {
			t.Errorf("got %d calls on gameLoop want at least %d", processor.processCalls, 60)
		}
	})

	t.Run("game loop stops after quit message", func(t *testing.T) {
		processor := stubWorldProcessor{}
		playerStore := make(PlayerStore)
		quitChan := make(chan struct{})
		endGameChan := make(chan struct{})
		actionQueue := &ActionQueue{actions: []Action{{input: MoveDown, player: "John"}}}
		stubEndGame := func() {
			close(endGameChan)
		}
		go gameLoop(&processor, playerStore, actionQueue, quitChan, stubEndGame)
		close(quitChan)

		select {
		case <-endGameChan:
			return
		case <-time.After(1 * time.Second):
			t.Errorf("stubEndGame wasn't called in time")
		}
	})
}

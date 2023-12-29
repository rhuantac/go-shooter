package main

import "testing"

func TestOperations(t *testing.T) {
	t.Run("actions are pushed to queue correctly", func(t *testing.T) {
		queue := ActionQueue{}
		for i := 0; i < 3; i ++ {
			queue.Push(Action{input: MoveDown, player: "John"})
		}

		if len(queue.actions) != 3 {
			t.Errorf("got %d queue actions want %d", len(queue.actions), 3)
		}

	})

	t.Run("queue is empty when actions are popped", func(t *testing.T) {
		queue := ActionQueue{}
		for i := 0; i < 3; i ++ {
			queue.Push(Action{input: MoveDown, player: "John"})
		}
		queue.PopAll()
		
		if len(queue.actions) != 0 {
			t.Errorf("got %d queue actions want %d", len(queue.actions), 3)
		}

	})
}
package main

import "testing"

func TestDemo(t *testing.T) {
	t.Run("snapshot consume works", func(t *testing.T) {
		cases := []struct {
			size int
			want int
		}{
			{50, 49},
			{1, 0},
			{0, 0},
		}

		for _, test := range cases {
			p := &Player{X: 0, Y: 0, states: make([]State, 0)}
			for i := 0; i < test.size; i++ {
				p.AddSnapshot(State{posX: 1, posY: 1, name: "John"})
			}

			p.RemoveSnapshot()

			if len(p.states) != test.want {
				t.Errorf("got %d snapshots want %d", len(p.states), test.want)
			}
		}

	})
}

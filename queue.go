package main

type ActionQueue struct {
	actions []Action
}

func (aq *ActionQueue) Push(action Action) {
	aq.actions = append(aq.actions, action)
}

func (aq *ActionQueue) PopAll() []Action {
	actionsBuffer := make([]Action, len(aq.actions))
	copy(actionsBuffer, aq.actions)
	aq.actions = aq.actions[:0]
	return actionsBuffer
}
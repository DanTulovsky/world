package world

import (
	"fmt"
	"sync"
)

var (
	actions = []action{"move", "skip"}
	wg      sync.WaitGroup
)

// Action describes what an exister can do on a given turn
type action string
type priority int32

type priorityAction struct {
	p priority
	a action
	f func()
}

// doActions dispatches an action to each peep.
func (w *World) doActions() {
	for _, l := range w.allLocations() {
		e := w.LocationExister(l)
		if e == nil {
			continue // empty spot
		}
		// get the best action for this Exister right now
		action := w.bestAction(e)
		// and run it
		action()
	}

}

// actionPriority returns a priority and the closure for the given action
// A higher priority will get executed first.
func (w *World) actionPriority(a action, e Exister, c chan priorityAction) {
	var p priority
	f := func() { return }

	switch a {
	case "move":
		f = func() {
			w.movePeep(e, allowMoves)
		}
		p = 1

	case "skip":
		p = 0

	default:

	}
	c <- priorityAction{p, a, f}
	close(c)
}

// bestAction returns a closure that executes the best action for a peep
func (w *World) bestAction(e Exister) func() {

	//possibleActions := make(map[priority]action)

	var channels []chan priorityAction

	for _, a := range actions {
		// wg.Add(1)
		c := make(chan priorityAction, 1)
		channels = append(channels, c)
		go w.actionPriority(a, e, c)
	}

	var bestAction func()
	var highestPriority priority

	// get results from all channels, keep only the one with max priority
	for _, c := range channels {
		Log("Processing channel!")
		var pa priorityAction
		for elem := range c {
			Log("--------", elem)
			pa = <-c
			Log("+++++++++", pa)
			if pa.p > highestPriority {
				highestPriority = pa.p
				bestAction = pa.f
			}

		}

	}

	Log(">>>>>>>>> ", bestAction)
	return bestAction
}

// movePeeps moves a single peep to the best location
func (w *World) movePeep(peep Exister, allowMoves bool) error {
	if !allowMoves {
		return fmt.Errorf("Moves not allowed by config.")
	}

	// Dead peeps don't move... for now.
	if !peep.IsAlive() {
		return fmt.Errorf("Dead peeps don't move!")
	}
	x, y, z := w.BestPeepMove(peep)

	if err := w.Move(peep, x, y, z); err != nil {
		return fmt.Errorf("Error moving peep: %v", err)
	}

	return nil
}

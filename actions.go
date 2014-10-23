package world

import (
	"fmt"
	"math/rand"
	"sync"
)

var (
	actions = []action{"move", "skip", "look"}
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

// bestAction returns a closure that executes the best action for a peep
func (w *World) bestAction(e Exister) func() {
	c := make(chan priorityAction, 5)

	for _, a := range actions {
		wg.Add(1)
		go w.actionPriority(a, e, c)
	}

	// this waits for all channels to return something and then closes the channel
	go func() {
		wg.Wait()
		close(c)
	}()

	var bestAction func()
	var highestPriority priority

	// get results from all goroutines, keep only the one with max priority
	// this ranges over c until the goroutine above closes c
	for pa := range c {
		if pa.p > highestPriority {
			highestPriority = pa.p
			bestAction = pa.f
		}

	}
	return bestAction
}

// lookAction returns a priority and the look action
// this trumps moveAction if the peep hasn't looked around recently
func (w *World) lookAction(e Exister) (priority, func()) {
	// Before moving, look around to see what's out there
	action := func() {
		w.LookAround(e)
	}
	pr := priority(10)

	if w.turn-e.LookTurn() < w.settings.PeepRememberTurns {
		pr = priority(1) // peeps still remembers the last time it looked
	}

	return pr, action
}

// moveAction returns the move action
func (w *World) moveAction(e Exister) (priority, func()) {
	action := func() {
		w.movePeep(e)
	}
	return 5, action
}

// skipAction returns the skip action
func (w *World) skipAction(e Exister) (priority, func()) {
	action := func() {
		return
	}
	return 0, action
}

// LookAround tells the Exister to look around and record what's around where
// Existers remember what they saw for 4 turns (configrable)
func (w *World) LookAround(e Exister) {
	e.SetNeighbors()
	e.SetLookTurn(w.turn)
}

// actionPriority returns a priority and the closure for the given action
// A higher priority will get executed first.
func (w *World) actionPriority(a action, e Exister, c chan priorityAction) {
	defer wg.Done()

	var p priority
	f := func() { return } // default is do nothing

	switch a {
	case "move":
		p, f = w.moveAction(e)

	case "look":
		p, f = w.lookAction(e)

	case "skip":
		p, f = w.skipAction(e)

	default:

	}
	c <- priorityAction{p, a, f}
}

// movePeeps moves a single peep to the best location
func (w *World) movePeep(peep Exister) error {
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

// BestPeepMove returns the most optimal move for a peep
func (w *World) BestPeepMove(e Exister) (int32, int32, int32) {
	// random for now
	var x, y, z int32
	// Peeps can move one square at a time in x, y direction.
	m := []int32{-1, 0, 1}
	x = m[rand.Intn(len(m))]
	y = m[rand.Intn(len(m))]
	return x, y, z
}

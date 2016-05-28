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
	var doneID []string // already processed

	for _, l := range w.allLocations() {

		e := w.LocationExister(l)

		if e == nil || ListContainsString(doneID, e.ID()) {
			continue // empty spot
		}

		// get the best action for this Exister right now
		action := w.bestAction(e)
		// and run it
		action()
		doneID = append(doneID, e.ID())
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
func (w *World) movePeep(e Exister) error {
	if !allowMoves {
		return fmt.Errorf("Moves not allowed by config.")
	}

	// Dead peeps don't move... for now.
	if !e.IsAlive() {
		return fmt.Errorf("Dead peeps don't move!")
	}

	// Look around first
	w.LookAround(e)

	x, y, z := w.BestPeepMove(e)

	if err := w.Move(e, x, y, z); err != nil {
		return fmt.Errorf("Error moving peep: %v", err)
	}

	return nil
}

// NextMoveToGetFromTo returns the x, y, z magnitude in order to move from src to dst
func (w *World) NextMoveToGetFromTo(src, dst Location) (x int32, y int32, z int32) {
	if src.SameAs(dst) {
		return 0, 0, 0
	}

	if dst.X > src.X {
		x = 1
	} else if dst.X < src.X {
		x = -1
	} else {
		x = 0
	}

	if dst.Y > src.Y {
		y = 1
	} else if dst.Y < src.Y {
		y = -1
	} else {
		y = 0
	}

	if dst.Z > src.Z {
		z = 1
	} else if dst.Z < src.Z {
		z = -1
	} else {
		z = 0
	}

	// check if the suggested square is busy and try alternatives
	if w.IsOccupiedLocation(Location{src.X + x, src.Y + y, src.Z + z}) {
		if !w.IsOccupiedLocation(Location{src.X, src.Y + y, src.Z + z}) &&
			!w.IsOutsideGrid(src.X, src.Y+y, src.Z+z) {
			return 0, y, z
		}
		if !w.IsOccupiedLocation(Location{src.X + x, src.Y, src.Z + z}) &&
			!w.IsOutsideGrid(src.X+x, src.Y, src.Z+z) {
			return x, 0, z
		}
		if !w.IsOccupiedLocation(Location{src.X + x, src.Y + y, src.Z}) &&
			!w.IsOutsideGrid(src.X+x, src.Y+y, src.Z) {
			return x, y, 0
		}
		// Random
		return w.randomMove()
	}
	return x, y, z
}

// NextMoveToGetAwayFrom returns the x, y, z magnitude in order to move away from loc while at current
func (w *World) NextMoveToGetAwayFrom(current, loc Location) (x int32, y int32, z int32) {

	if loc.X >= 0 {
		x = -1
	} else {
		x = 1
	}

	if loc.Y >= 0 {
		y = -1
	} else {
		y = 1
	}

	// check if the suggested square is busy and try alternatives
	if w.IsOccupiedLocation(Location{current.X + x, current.Y + y, current.Z + z}) ||
		w.IsOutsideGrid(current.X+x, current.Y+y, current.Z+z) {
		if !w.IsOccupiedLocation(Location{current.X, current.Y + y, current.Z + z}) &&
			!w.IsOutsideGrid(current.X, current.Y+y, current.Z+z) {
			return 0, y, z
		}
		if !w.IsOccupiedLocation(Location{current.X + x, current.Y, current.Z + z}) &&
			!w.IsOutsideGrid(current.X+x, current.Y, current.Z+z) {
			return x, 0, z
		}
		if !w.IsOccupiedLocation(Location{current.X + x, current.Y + y, current.Z}) &&
			!w.IsOutsideGrid(current.X+x, current.Y+y, current.Z) {
			return x, y, 0
		}
		// Random
		return w.randomMove()
	}
	return x, y, z
}

// BestPeepMove returns the most optimal move for a peep
// x, y and z are magnitudes, not coordinates.
func (w *World) BestPeepMove(e Exister) (x int32, y int32, z int32) {
	// Suggest move based on neighbors around
	neighbors := e.NeighborsFromLook()

	for _, n := range neighbors {
		// Move towards same gender if have not yet spawned and are both of spawn age
		if n.Gender() == e.Gender() {
			if w.turn - n.SpawnTurn() < w.settings.PeepSpawnInterval {
				Log("too recent spawn", n.SpawnTurn())
				continue // spawned too recently
			}

			if e.MetPeep(n) {
				continue // already met
			}

			if w.OfSpawnAge(e) && w.OfSpawnAge(n) {
				return w.NextMoveToGetFromTo(e.Location(), n.Location())
			}
		} else { // different genders
			// Move towards different gender peep
			return w.NextMoveToGetFromTo(e.Location(), n.Location())
		}
	}

	return w.randomMove()

	// Possibilities
	// No interesting neighbors around
	if w.OfSpawnAge(e) {
		// Move towards base
		return w.NextMoveToGetFromTo(e.Location(), e.Homebase())
	} else {
		// Move away from base
		return w.NextMoveToGetAwayFrom(e.Location(), e.Homebase())
	}

}

func (w *World) randomMove() (x, y, z int32) {
	m := []int32{-1, 0, 1}
	return m[rand.Intn(len(m))], m[rand.Intn(len(m))], z
}

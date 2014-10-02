package world

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var (
// random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Log(txt ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", txt)
}

// World describes the world state
type World struct {
	peeps      []*Peep // citizens
	name       string
	settings   Settings
	turn       int64              // the current turn
	eventQueue chan termbox.Event // for catching user input
	grid       *Grid              // Map of coordinates to occupant
}

func NewWorld(name string, settings Settings, eventQueue chan termbox.Event) *World {
	rand.Seed(time.Now().UnixNano())
	return &World{
		name:       name,
		settings:   settings,
		eventQueue: eventQueue, // keyboard input events
		grid: &Grid{
			size:    settings.Size,
			objects: NewDmap(), // empty grid
		},
	}
}

// NextTurn advances the world to the next turn.
func (w *World) NextTurn() error {
	// Check if we should exit
	select {
	case ev := <-w.eventQueue:
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			return errors.New("Exiting...")
		}
	default:
		// Redraw screen
		w.Draw()

		w.turn++

		// Move peeps around
		w.MovePeeps()

		// New peep might be born
		if err := w.randomPeep(); err != nil {
			Log(err)
		}

		// Age existing peeps
		for _, peep := range w.peeps {
			if !peep.IsAlive() {
				continue
			}
			peep.AgeOrDie(w.settings.MaxAge, w.settings.RandomDeath)
		}

	}
	return nil
}

// BestPeepMove returns the most optimal move for a peep
func (w *World) BestPeepMove(peep *Peep) (int32, int32, int32) {
	// random for now
	var x, y, z int32
	// Peeps can move one square at a time in x, y direction.
	m := []int32{-1, 0, 1}
	x = m[rand.Intn(len(m))]
	y = m[rand.Intn(len(m))]
	return x, y, z
}

// MovePeeps moves peeps around every turn
func (w *World) MovePeeps() {
	for _, peep := range w.peeps {
		// Dead peeps don't move... for now.
		if !peep.IsAlive() {
			continue
		}
		x, y, z := w.BestPeepMove(peep)

		//Log(fmt.Sprintf("Moving %v (%v): (%v, %v, %v)", peep.ID(), peep.Location(), x, y, z))
		if err := w.Move(peep, x, y, z); err != nil {
			Log(err)
		}
	}
}

// randomPeep creates a new peep at random
// randomness controlled by world.settings.NewPeepModifier
// As the world grows, probability of this event goes towards 0
// Subject to world.settings.MaxPeeps
func (w *World) randomPeep() error {
	// MaxPeeps already, short circuit here.
	if w.AlivePeeps() >= w.settings.MaxPeeps || w.AlivePeeps() >= w.settings.NewPeepMax {
		return fmt.Errorf("cannot create new peep, MaxPeeps already present")
	}

	// Something at origin
	e := w.grid.objects.GetByLocation(Location{0, 0, 0})

	if e != nil && w.IsAlive(e.ID()) {
		return fmt.Errorf("cannot crate new peep, origin taken by: %v", e.ID())
	}

	probability := w.settings.NewPeep - (float64(w.AlivePeeps()) / w.settings.NewPeepModifier)
	if rand.Float64() < probability {
		w.NewPeep("", Location{})
	}
	return nil
}

// AlivePeeps returns the number of alive peeps
func (world *World) AlivePeeps() int64 {
	var peeps int64
	for _, p := range world.peeps {
		if p.IsAlive() {
			peeps++
		}
	}
	return peeps
}

// PeepGenders returns a count of all peep genders
func (world *World) PeepGenders() map[PeepGender]int64 {
	genders := make(map[PeepGender]int64)
	for _, p := range world.peeps {
		if p.IsAlive() {
			genders[p.Gender()]++
		}
	}
	return genders
}

// PeepMaxAge returns the max age of all peeps
func (w *World) PeepMaxAge() PeepAge {
	var max PeepAge
	for _, p := range w.peeps {
		if p.Age() > max && p.IsAlive() {
			max = p.Age()
		}
	}
	return max
}

// PeepMinAge returns the min age of all peeps
func (w *World) PeepMinAge() PeepAge {
	if w.AlivePeeps() == 0 {
		return 0
	}
	min := w.settings.MaxAge
	for _, p := range w.peeps {
		if p.Age() < min && p.IsAlive() {
			min = p.Age()
		}
	}
	return min
}

// PeepAvgAge returns the average age of all peeps
func (w *World) PeepAvgAge() PeepAge {
	var sum PeepAge
	var alive PeepAge
	for _, p := range w.peeps {
		if p.IsAlive() {
			sum += p.Age()
			alive++
		}
	}
	if alive > 0 {
		return sum / alive
	}
	return 0
}

// Run runs the world.
func (w *World) Run() {
	Log("Starting world...")
}

// Pause pauses the world.
func (w *World) Pause() {

}

// String prints world information.
func (world *World) Show() {
	io := os.Stderr
	fmt.Fprintf(io, "%v\n", strings.Repeat("-", 80))
	fmt.Fprintf(io, "Name: %v\n", world.name)
	fmt.Fprintf(io, "Turn: %v\n", world.turn)
	fmt.Fprintf(io, "Peeps: %v/%v\n", world.AlivePeeps(), world.settings.MaxPeeps)
	fmt.Fprintf(io, "Absolute MaxAge: %v\n", world.settings.MaxAge)
	fmt.Fprintf(io, "Peep Max/Avg/Min Age: %v/%v/%v\n", world.PeepMaxAge(), world.PeepAvgAge(), world.PeepMinAge())
	fmt.Fprintf(io, "Genders: %v\n", world.PeepGenders())

	Log("World GRID:")
	Log(strings.Repeat("*", 40))
	for _, peep := range world.peeps {
		if peep.IsAlive() {
			Log("%%%%", peep.ID(), peep.Location())
		}
	}
	Log(strings.Repeat("*", 40))
}

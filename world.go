package world

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
	termbox "github.com/nsf/termbox-go"
)

var (
	// random = rand.New(rand.NewSource(time.Now().UnixNano()))
	allowMoves = true // for testing, turns off random moves.
)

func Log(txt ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", txt)
}

type neighborViewDistanceCache struct {
	l Location
	v int32
}

// World describes the world state
type World struct {
	peeps             []*Peep // citizens
	name              string
	settings          Settings
	turn              Turn               // the current turn
	eventQueue        chan termbox.Event // for catching user input
	grid              *Grid              // Map of coordinates to occupant
	stats             *stats
	locationNeighbors map[neighborViewDistanceCache][]Location // cache of location/view distance -> list of neighbor locations
	okToAdvance       bool                                     // for debugging
	debug             bool
	homebase          map[PeepGender]Location
}

type Turn int64

// ListContains returns true if Location is in the list.
func ListContains(list []Location, loc Location) bool {
	for _, l := range list {
		if l.SameAs(loc) {
			return true
		}
	}
	return false
}

// ListContainsString returns true if a string is in the list
func ListContainsString(list []string, s string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}
	return false
}

func NewWorld(name string, settings Settings, eventQueue chan termbox.Event, debug bool) *World {
	rand.Seed(time.Now().UnixNano())
	return &World{
		name:       name,
		settings:   settings,
		eventQueue: eventQueue, // keyboard input events
		grid: &Grid{
			size:    settings.Size,
			objects: NewDmap(), // empty grid
		},
		stats:             newStats(),
		locationNeighbors: make(map[neighborViewDistanceCache][]Location),
		debug:             debug,
		homebase:          make(map[PeepGender]Location),
	}
}

// handleOvercrowding handles the cases when a peep is completely surrounded
func (w *World) handleOvercrowding(p *Peep) {
	// Get all neighboring locations
	neighborLocations := w.LocationNeighbors(p.Location(), 1)

	var genderCount = make(map[PeepGender]int)
	for _, l := range neighborLocations {
		if e := w.LocationExister(l); e != nil {
			genderCount[e.Gender()]++
		}
	}

	// Check if surrounded and kill if settings say so
	if w.settings.KillIfSurroundByOther {
		var otherNeighbors int

		for gender, count := range genderCount {
			if gender != p.Gender() {
				otherNeighbors += count
			}
		}
		if len(neighborLocations) == otherNeighbors {
			p.Die(w.turn)
		}

	}

	if w.settings.KillIfSurroundedBySame {
		// If all locations around are take up by same gender peeps
		if len(neighborLocations) == genderCount[p.Gender()] {
			p.Die(w.turn)
		}
	}

	if w.settings.KillIfSurrounded {
		var allNeighbors int
		for _, count := range genderCount {
			allNeighbors += count
		}

		if len(neighborLocations) == allNeighbors {
			p.Die(w.turn)
		}
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
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeySpace {
			w.Show(os.Stderr)
		}
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlS {
			w.ShowSettings(os.Stderr)
		}
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
			w.okToAdvance = true
		}
	default:
		if w.debug {
			if !w.okToAdvance {
				return nil
			} else {
				Log("Advancing to next turn...")
			}

		}

		// Update stats
		w.stats.peepsAlive.Update(w.AlivePeepCount())
		w.stats.peepsDead.Update(w.DeadPeepCount())

		// Redraw screen
		w.Draw()

		w.turn++

		// Peep actions
		w.doActions()

		// New peep might be born
		if err := w.randomPeep(); err != nil {
			//Log(err)
		}

		// Age and/or kill existing peeps
		for _, peep := range w.peeps {
			if !peep.IsAlive() {
				continue
			}
			age, err := peep.AgeOrDie(w.settings.MaxAge, w.settings.RandomDeath, w.turn)
			if err != nil {
				w.stats.ages.Update(int64(age))
			}
			w.handleOvercrowding(peep)

		}

		if w.debug {
			w.okToAdvance = false
			w.Show(os.Stderr)
			w.ShowGrid(os.Stderr)
		}

	}
	return nil
}

// randomPeep creates a new peep at random
// randomness controlled by world.settings.NewPeepModifier
// As the world grows, probability of this event goes towards 0
// Subject to world.settings.MaxPeeps
func (w *World) randomPeep() error {
	if w.AlivePeepCount() >= w.settings.NewPeepMax {
		return fmt.Errorf("Too many peeps (%v) for random spawn.", w.AlivePeepCount())
	}
	probability := w.settings.NewPeep - (float64(w.AlivePeepCount()) / w.settings.NewPeepModifier)
	if rand.Float64() < probability {
		w.NewPeep("", Location{})
	}
	return nil
}

// AlivePeepCount returns the number of alive peeps
func (world *World) AlivePeepCount() int64 {
	var peeps int64
	for _, p := range world.peeps {
		if p.IsAlive() {
			peeps++
		}
	}
	return peeps
}

// DeadPeepCount returns the number of dead peeps
func (world *World) DeadPeepCount() int64 {
	var peeps int64
	for _, p := range world.peeps {
		if !p.IsAlive() {
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
	if w.AlivePeepCount() == 0 {
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

func (w *World) runWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", w.HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":6001", nil)
}

// Run runs the world.
func (w *World) Run() {
	Log("Starting world...")
	go w.runWebServer()
}

// Pause pauses the world.
func (w *World) Pause() {

}

// ShowGrid prints the grid and its occupants
func (world *World) ShowGrid(w io.Writer) {
	fmt.Fprintf(w, "World GRID:\n")
	fmt.Fprintf(w, "%v\n", strings.Repeat("*", 40))
	for _, peep := range world.peeps {
		if peep.IsAlive() {
			fmt.Fprintf(w, "%v\n", peep.String())
		}
	}
	fmt.Fprintf(w, "%v\n", strings.Repeat("*", 40))
}

// String prints world information.
func (world *World) Show(w io.Writer) {
	fmt.Fprintf(w, "%v\n", strings.Repeat("-", 80))
	fmt.Fprintf(w, "Name: %v\n", world.name)
	fmt.Fprintf(w, "Turn: %v\n", world.turn)
	fmt.Fprintf(w, "Peeps Alive/Dead/MaxAlive: %v/%v/%v\n", world.AlivePeepCount(), world.DeadPeepCount(), world.settings.MaxPeeps)
	fmt.Fprintf(w, "Peep Max/Avg/Min Age: %v/%v/%v\n", world.PeepMaxAge(), world.PeepAvgAge(), world.PeepMinAge())
	fmt.Fprintf(w, "Genders: %v\n", world.PeepGenders())

}

// String prints world information.
func (world *World) ShowSettings(w io.Writer) {
	fmt.Fprintf(w, "%v\n", strings.Repeat("-", len("Settings")))
	fmt.Fprintf(w, "Settings\n")
	fmt.Fprintf(w, "%v\n", strings.Repeat("-", len("Settings")))

	s := reflect.ValueOf(&world.settings).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Fprintf(w, "%s = %v\n", typeOfT.Field(i).Name, f.Interface())
	}
}

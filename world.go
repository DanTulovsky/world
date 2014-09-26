package world

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
// random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Log(txt ...interface{}) {
	fmt.Printf("%v\n", txt)
}

// World describes the world state
type World struct {
	peeps    []*Peep // citizens
	name     string
	settings Settings
	turn     int64 // the current turn
}

func NewWorld(name string, settings Settings) *World {
	rand.Seed(time.Now().UnixNano())
	return &World{
		name:     name,
		settings: settings,
	}
}

// NextTurn advances the world to the next turn.
func (world *World) NextTurn() {
	world.turn++

	// New peep might be born
	world.newPeep()

	// Age existing peeps
	for _, peep := range world.peeps {
		peep.AgeOrDie(world.settings.MaxAge, world.settings.RandomDeath)
	}
}

// newPeep creates a new peep at random
// randomness controlled by world.settings.newpeep
// As the world grows, probability of this event goes towards 0
func (world *World) newPeep() {
	probability := world.settings.NewPeep - (float64(world.AlivePeeps()) / world.settings.NewPeepModifier)
	if rand.Float64() < probability {
		world.peeps = append(world.peeps, NewPeep())
	}
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
func (world *World) PeepMaxAge() PeepAge {
	var max PeepAge
	for _, p := range world.peeps {
		if p.Age() > max && p.IsAlive() {
			max = p.Age()
		}
	}
	return max
}

// PeepMinAge returns the min age of all peeps
func (world *World) PeepMinAge() PeepAge {
	min := world.settings.MaxAge
	for _, p := range world.peeps {
		if p.Age() < min && p.IsAlive() {
			min = p.Age()
		}
	}
	return min
}

// PeepAvgAge returns the average age of all peeps
func (world *World) PeepAvgAge() PeepAge {
	var sum PeepAge
	var alive PeepAge
	for _, p := range world.peeps {
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
func (world *World) Run() {
	Log("Starting world...")
}

// Pause pauses the world.
func (world *World) Pause() {

}

// String prints world information.
func (world *World) Show() {
	fmt.Printf("%v\n", strings.Repeat("-", 80))
	fmt.Printf("Name: %v\n", world.name)
	fmt.Printf("Turn: %v\n", world.turn)
	fmt.Printf("Peeps: %v\n", world.AlivePeeps())
	fmt.Printf("Peep Max/Avg/Min Age: %v/%v/%v\n", world.PeepMaxAge(), world.PeepAvgAge(), world.PeepMinAge())
	fmt.Printf("Genders: %v\n", world.PeepGenders())

	//for _, peep := range world.peeps {
	//	if peep.IsAlive() {
	//		fmt.Printf("%v\n", peep)
	//	}
	//}
}

package world

import (
	"fmt"
	"math"
	"math/rand"

	"code.google.com/p/go-uuid/uuid"
)

var (
	genders  = []PeepGender{"blue", "red", "green", "yellow"}
	homebase = make(map[PeepGender]Location)
)

type PeepAge int64
type PeepGender string

type Peep struct {
	id         string // unique id
	age        PeepAge
	isalive    bool
	gender     PeepGender
	deadAtTurn int64             // World turn when the peep died
	met        map[Exister]int64 // records all other existers and the turn met
}

func (w *World) Genders() []PeepGender {
	return genders
}

func (w *World) SetHomebase(gender PeepGender, loc Location) {
	homebase[gender] = loc
}

// NewPeep creates and returns a new peep
func (w *World) NewPeep(gender PeepGender, location Location) (*Peep, error) {
	// MaxPeeps already
	if w.AlivePeeps() >= w.settings.MaxPeeps {
		return nil, fmt.Errorf("cannot create new peep, MaxPeeps already present")
	}

	if gender == "" {
		gender = genders[rand.Intn(len(genders))]
	}
	peep := &Peep{
		id:      uuid.New(),
		isalive: true,
		gender:  gender,
		met:     make(map[Exister]int64),
	}
	// If no specific location set, pick one based on gender
	if location.SameAs(Location{}) {
		location = w.SpawnPoint(peep)
	}

	// Check if spawn point is busy.
	e := w.grid.objects.GetByLocation(location)
	if e != nil && e.IsAlive() {
		return nil, fmt.Errorf("cannot crate new peep, origin taken by: %v", e.ID())
	}

	w.peeps = append(w.peeps, peep)
	w.UpdateGrid(peep, location, location)
	return peep, nil
}

func (peep *Peep) ID() string {
	return peep.id
}

func (peep *Peep) String() string {
	return fmt.Sprintf("%v age:%v gender:%v", peep.id, peep.age, peep.gender)
}

// Homebase returns the homebase location given a peep
func (peep *Peep) Homebase() Location {
	return homebase[peep.Gender()]
}

// IsAlive returns True of peep is alive.
func (peep *Peep) IsAlive() bool {
	return peep.isalive
}

// Age returns peep's Age
func (peep *Peep) Age() PeepAge {
	return peep.age
}

func (peep *Peep) DeadAtTurn() int64 {
	return peep.deadAtTurn
}

// AddAge increases the age of the peep by 1
func (peep *Peep) AddAge() {
	peep.age++
}

// Meet records a meeting between peep and other
// This records both sides
func (peep *Peep) Meet(other Exister, turn int64) {
	peep.met[other] = turn
}

// Met returns all the existers this one met, and the turn.
func (peep *Peep) Met() map[Exister]int64 {
	return peep.met
}

// MetPeep returns true if two peeps have met
func (peep *Peep) MetPeep(other Exister) bool {
	if _, ok := peep.met[other]; ok {
		return true
	}
	return false
}

// Die kills the peep
func (peep *Peep) Die(turn int64) {
	// Log("Peep: ", peep.ID(), " died!")
	peep.isalive = false
	peep.deadAtTurn = turn
}

// Gender returns the peep's gender
func (peep *Peep) Gender() PeepGender {
	return peep.gender
}

// AgeOrDie ages a peep or kills him
// based on age and probability
func (peep *Peep) AgeOrDie(maxage PeepAge, randomdeath float64, turn int64) {
	if peep.age >= maxage {
		peep.Die(turn)
		return
	}
	// Older peeps have more chances to die
	if rand.Float64() < randomdeath+(math.Log10(float64(peep.age))/float64(maxage/1)) {
		peep.Die(turn)
		return
	}
	peep.AddAge()
}

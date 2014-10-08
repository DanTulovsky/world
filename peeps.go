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
	id      string // unique id
	age     PeepAge
	isalive bool
	gender  PeepGender
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

// AddAge increases the age of the peep by 1
func (peep *Peep) AddAge() {
	peep.age++
}

// Die kills the peep
func (peep *Peep) Die() {
	// Log("Peep: ", peep.ID(), " died!")
	peep.isalive = false
}

// Gender returns the peep's gender
func (peep *Peep) Gender() PeepGender {
	return peep.gender
}

// AgeOrDie ages a peep or kills him
// based on age and probability
func (peep *Peep) AgeOrDie(maxage PeepAge, randomdeath float64) {
	if peep.age >= maxage {
		peep.Die()
		return
	}
	// Older peeps have more chances to die
	if rand.Float64() < randomdeath+(math.Log10(float64(peep.age))/float64(maxage/1)) {
		peep.Die()
		return
	}
	peep.AddAge()
}

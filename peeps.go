package world

import (
	"fmt"
	"math"
	"math/rand"

	"code.google.com/p/go-uuid/uuid"
)

var (
	genders = []PeepGender{"blue", "red", "green", "yellow"}
)

type PeepAge int64
type PeepGender string

type Peep struct {
	id         string // unique id
	age        PeepAge
	isalive    bool
	gender     PeepGender
	deadAtTurn Turn                 // World turn when the peep died
	met        map[Exister]Turn     // records all other existers and the turn met
	lookTurn   Turn                 // the turn when this peep looked around
	world      *World               // reference to world
	neighbors  map[Location]Exister // neighbors at time of last lookup
	spawnTurn  Turn                 // the turn of last spawn
}

func (w *World) Genders() []PeepGender {
	return genders[0:w.settings.MaxGenders]
}

func (w *World) SetHomebase(gender PeepGender, loc Location) {
	w.homebase[gender] = loc
}

// NewPeep creates and returns a new peep
func (w *World) NewPeep(gender PeepGender, location Location) (*Peep, error) {
	// MaxPeeps already
	if w.AlivePeepCount() >= w.settings.MaxPeeps {
		return nil, fmt.Errorf("cannot create new peep, MaxPeeps already present")
	}

	if gender == "" {
		gender = genders[rand.Intn(len(w.Genders()))]
	}
	peep := &Peep{
		id:        uuid.New(),
		isalive:   true,
		gender:    gender,
		met:       make(map[Exister]Turn),
		world:     w,
		neighbors: make(map[Location]Exister),
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

	w.UpdateGrid(peep, location, location)
	return peep, nil
}

// Location returns this peep's location
func (p *Peep) Location() Location {
	l, _ := p.world.ExisterLocation(p)
	return l
}

// NeighborsFromLook returns the map of neighbors at time of last lookup
func (p *Peep) NeighborsFromLook() map[Location]Exister {
	return p.neighbors
}

// SetNeighbors sets exister's neighbors right now
// All neighbors in the radius of world.settings.PeepViewDistance are returned
func (p *Peep) SetNeighbors() {

	locations := p.world.LocationNeighbors(p.Location(), p.world.settings.PeepViewDistance)

	for _, l := range locations {
		if p.Location().SameAs(l) {
			continue
		}
		e := p.world.LocationExister(l)
		if e != nil && e.IsAlive() { // don't care about dead existers
			p.neighbors[l] = e
		}
	}
}

// SpawnTurn returns the last time the peep spawned
func (p *Peep) SpawnTurn() Turn {
	return p.spawnTurn
}

// SetSpawnTurn sets the spawn turn
func (p *Peep) SetSpawnTurn(t Turn) {
	p.spawnTurn = t
}

// World returns pointer to the world the exister is in
func (p *Peep) World() *World {
	return p.world
}

// LookTurn returns the last turn that the peep looked around
func (p *Peep) LookTurn() Turn {
	return p.lookTurn
}

// SetLookTurn sets the last turn this peep looked around
func (p *Peep) SetLookTurn(t Turn) {
	p.lookTurn = t
}

func (peep *Peep) ID() string {
	return peep.id
}

func (peep *Peep) String() string {
	return fmt.Sprintf("%v age:%v gender:%v location:%v", peep.ID(), peep.Age(), peep.Gender(), peep.Location())
}

// Homebase returns the homebase location given a peep
func (peep *Peep) Homebase() Location {
	return peep.world.homebase[peep.Gender()]
}

// IsAlive returns True of peep is alive.
func (peep *Peep) IsAlive() bool {
	return peep.isalive
}

// Age returns peep's Age
func (peep *Peep) Age() PeepAge {
	return peep.age
}

func (peep *Peep) DeadAtTurn() Turn {
	return peep.deadAtTurn
}

// AddAge increases the age of the peep by 1
func (peep *Peep) AddAge() {
	peep.age++
}

// Meet records a meeting between peep and other
// This records both sides
func (peep *Peep) Meet(other Exister, turn Turn) {
	peep.met[other] = turn
}

// Met returns all the existers this one met, and the turn.
func (peep *Peep) Met() map[Exister]Turn {
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
func (peep *Peep) Die(turn Turn) {
	// Log("Peep: ", peep.ID(), " died!")
	peep.isalive = false
	peep.deadAtTurn = turn
}

// Gender returns the peep's gender
func (peep *Peep) Gender() PeepGender {
	return peep.gender
}

// AgeOrDie ages a peep or kills him based on age and probability
// An error is return on death
func (peep *Peep) AgeOrDie(maxage PeepAge, randomdeath float64, turn Turn) (PeepAge, error) {
	if peep.age >= maxage {
		peep.Die(turn)
		return peep.Age(), fmt.Errorf("Peep died, too old...")
	}
	// Older peeps have more chances to die
	if randomdeath > 0 && rand.Float64() < randomdeath+(math.Log10(float64(peep.age))/float64(maxage/1)) {
		peep.Die(turn)
		return peep.Age(), fmt.Errorf("Peep died, randomness sucks...")
	}
	peep.AddAge()
	return peep.Age(), nil
}

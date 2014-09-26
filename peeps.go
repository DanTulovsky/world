package world

import (
	"fmt"
	"math"
	"math/rand"
)

var (
	genders = []PeepGender{"blue", "red"}
)

type PeepAge int64
type PeepGender string

type Peep struct {
	id       float64 // unique id
	age      PeepAge
	isalive  bool
	gender   PeepGender
	location Location
}

// NewPeep returns a new peep
func NewPeep() *Peep {
	return &Peep{
		id:      rand.Float64(),
		isalive: true,
		gender:  genders[rand.Intn(len(genders))],
	}
}

func (peep *Peep) String() string {
	return fmt.Sprintf("%v age:%v gender:%v", peep.id, peep.age, peep.gender)
}

// Location returns the peep's location
func (peep *Peep) Location() Location {
	return peep.location
}

// MoveX implements Location interface
func (peep *Peep) MoveX(steps int32) {
	peep.location.x += steps
}

// MoveY implements Location interface
func (peep *Peep) MoveY(steps int32) {
	peep.location.y += steps
}

// MoveZ implements Location interface
func (peep *Peep) MoveZ(steps int32) {
	peep.location.z += steps
}

// SetLocation sets a peep's location
//func (peep *Peep) Move(direction Direction) {
//	peep.location.X += direction.X
//	peep.location.Y += direction.Y
//	peep.location.Z += direction.Z
//}

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
	if rand.Float64() < randomdeath+(math.Log10(float64(peep.age))/float64(maxage/3)) {
		peep.Die()
		return
	}
	peep.AddAge()
}

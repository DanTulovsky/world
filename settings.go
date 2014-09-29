package world

import "time"

type Settings struct {
	MaxAge   PeepAge // cannot live beyond this age
	MaxPeeps int64   // Absolute max peeps
	NewPeep  float64 // Chances a new peep is born [0-1]

	// The lower this number, the less chance a new peep will show up as the population grows
	// At 1, new random peeps will almost never show up
	NewPeepModifier  float64
	NewPeepMax       int64         // When this many peeps exist, no new peeps are spawned from origin
	RandomDeath      float64       // chances of a random death
	Size             *Size         // world size, one line is used as the border around
	SpawnProbability float64       // chances of two peeps that meet spawning a new one
	TurnTime         time.Duration // How fast is each turn?
}

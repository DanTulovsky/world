package world

type Settings struct {
	// Chances of events
	NewPeep     float64 // new peep is born
	MaxAge      PeepAge // cannot live beyond this age
	RandomDeath float64 // chances of a random death

	// The lower this number, the less chance a new peep will show up as the population grows
	// At 1, new random peeps will almost never show up
	NewPeepModifier float64
}

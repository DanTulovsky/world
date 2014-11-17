package world

import "time"

type Settings struct {
	MaxAge   PeepAge // cannot live beyond this age
	MaxPeeps int64   // Absolute max peeps
	NewPeep  float64 // Chances a new peep is born [0-1]

	// The lower this number, the less chance a new peep will show up as the population grows
	// At 1, new random peeps will almost never show up
	NewPeepModifier        float64
	NewPeepMax             int64         // When this many peeps exist, no new peeps are spawned from origin
	RandomDeath            float64       // chances of a random death
	Size                   *Size         // world size, one line is used as the border around
	SpawnAge               PeepAge       // Minimum age to spaw
	SpawnProbability       float64       // chances of two peeps that meet spawning a new one
	TurnTime               time.Duration // How fast is each turn?
	YoungHightlightAge     PeepAge       // Up to this age, peeps are highlighted in the GUI
	PeepRememberTurns      Turn          // How many turns peeps remember their surroundings for
	PeepViewDistance       int32         // how far they can see
	PeepSpawnInterval      Turn          // How many turns to wait after a spawn before can spawn again
	KillIfSurroundByOther  bool          // If surrounded completely by other genders, die
	KillIfSurroundedBySame bool          // If surrouned completelt by same genders, die
	KillIfSurrounded       bool          // If surrounded completely, die
	MaxGenders             int           // Maximux different genders.  1-4
}

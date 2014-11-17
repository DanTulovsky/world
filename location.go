package world

import (
	"fmt"
	"math"
	"math/rand"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

// Grid holds information about what object is in what cell in the world.
type Grid struct {
	size *Size

	// Each object is represented by a unique string
	objects *dmap
}

// Location specifies one coordinate in the world.
type Location struct {
	X int32
	Y int32
	Z int32
}

func (l Location) String() string {
	return fmt.Sprintf("(%v, %v, %v)", l.X, l.Y, l.Z)
}

// SameAs returns true if two locations are the same
func (l Location) SameAs(other Location) bool {

	if l.X == other.X && l.Y == other.Y && l.Z == other.Z {
		return true
	}
	return false

}

// NewLocation returns a new location at origin
func NewLocation() Location {
	return Location{0, 0, 0}
}

// NewLocationXYZ returns a new location at x,y,z
func NewLocationXYZ(x, y, z int32) Location {
	return Location{x, y, z}
}

// Size specifies the world size
// termbox starts at (0, 0) upper-left.
type Size struct {
	MaxX int32
	MaxY int32
	MaxZ int32

	MinX int32
	MinY int32
	MinZ int32
}

// Exister defines an object that exists in the grid
type Exister interface {
	ID() string
	Gender() PeepGender
	Age() PeepAge
	IsAlive() bool
	Homebase() Location
	DeadAtTurn() Turn
	Met() map[Exister]Turn // Map of exister to turn when met
	Meet(Exister, Turn)
	MetPeep(Exister) bool                    // Whether the two have met
	SetLookTurn(Turn)                        // sets the turn the exister looked around
	LookTurn() Turn                          // last time exister looked around
	World() *World                           // returns pointer to the World this exister inhabits
	SetNeighbors()                           // sets the neighbors around the exister on this turn
	NeighborsFromLook() map[Location]Exister // gets the neighbors of the exister (from the last look time)
	Location() Location                      // location of the exister on the map
	SpawnTurn() Turn                         // last time exister spawned
	SetSpawnTurn(Turn)                       // sets the spawn turn
}

// MaxX returns the max X value of the grid that can be occupied
func (w *World) MaxX() int32 {
	return w.settings.Size.MaxX - 1
}

// MaxY returns the max X value of the grid that can be occupied
func (w *World) MaxY() int32 {
	return w.settings.Size.MaxY - 1
}

// MaxZ returns the max Z value of the grid that can be occupied
func (w *World) MaxZ() int32 {
	return w.settings.Size.MaxZ
}

// MinX returns the min X value of the grid that can be occupied
func (w *World) MinX() int32 {
	return w.settings.Size.MinX + 1
}

// MinY returns the min Y value of the grid that can be occupied
func (w *World) MinY() int32 {
	return w.settings.Size.MinY + 1
}

// MinZ returns the min Z value of the grid that can be occupied
func (w *World) MinZ() int32 {
	return w.settings.Size.MinZ
}

// SpawnLocations returns all locations available for spawning
func (w *World) SpawnLocations() []Location {
	l := []Location{}
	// top left
	l = append(l, Location{w.MinX(), w.MaxY(), 0})
	// top right
	l = append(l, Location{w.MaxX(), w.MaxY(), 0})

	// bottom left
	l = append(l, Location{w.MinX(), w.MinY(), 0})
	// bottom right
	l = append(l, Location{w.MaxX(), w.MinY(), 0})

	return l
}

// LocationNeighbors returns all neighboring locations to the given one within the viewDistance
func (w *World) LocationNeighbors(l Location, viewDistance int32) []Location {
	// check cache first
	if neighbors, ok := w.locationNeighbors[neighborViewDistanceCache{l, viewDistance}]; ok {
		return neighbors
	}

	neighbors := []Location{}

	for x := -viewDistance; x <= viewDistance; x++ {
		for y := -viewDistance; y <= viewDistance; y++ {
			newLoc := NewLocationXYZ(l.X+x, l.Y+y, l.Z)
			if newLoc.SameAs(l) {
				continue // skip our own location
			}
			if !w.IsOutsideGrid(newLoc.X, newLoc.Y, newLoc.Z) {
				neighbors = append(neighbors, newLoc)
			}
		}
	}
	// Update cache
	w.locationNeighbors[neighborViewDistanceCache{l, viewDistance}] = neighbors
	return neighbors
}

// allLocations returns a list of all available locations
// may need to support re-sizing in the future
func (w *World) allLocations() []Location {
	all := make([]Location, 1)

	var z int32 // World is flat for now.
	for x := w.MinX(); x <= w.MaxX(); x++ {
		for y := w.MinY(); y <= w.MaxY(); y++ {
			all = append(all, NewLocationXYZ(x, y, z))
		}
	}
	return all
}

// FindAnyEmptyLocation returns the first empty location it finds.
func (w *World) FindAnyEmptyLocation() (Location, error) {
	var z int32 // World is flat for now.
	for x := w.MinX(); x <= w.MaxX(); x++ {
		for y := w.MinY(); y <= w.MaxY(); y++ {
			loc := NewLocationXYZ(x, y, z)
			if !w.IsOccupiedLocation(loc) {
				return loc, nil
			}
		}
	}
	return Location{}, fmt.Errorf("Unable to find empty location!")
}

// FindEmptyLocation returns an empty location next to one of the provided locations or an error if not able to find one
func (w *World) FindEmptyLocation(locations ...Location) (Location, error) {
	for _, l := range locations {
		neighbors := w.LocationNeighbors(l, 1)
		for _, n := range neighbors {
			if !w.IsOccupiedLocation(n) {
				return n, nil
			}
		}
	}
	return Location{}, fmt.Errorf("No available locations next to %v", locations)
}

// OfSpawnAge returns true of Exister is old enough to spawn
func (w *World) OfSpawnAge(e Exister) bool {
	return e.Age() >= w.settings.SpawnAge
}

// SameGenderSpawn makes a new peep next to one of the provided peeps, if they are of the same gender
func (w *World) SameGenderSpawn(left, right Exister) error {
	if left.Gender() != right.Gender() {
		return fmt.Errorf("Different genders don't spawn!")
	}
	if !w.OfSpawnAge(left) || !w.OfSpawnAge(right) {
		return fmt.Errorf("Both must be of spawn age!")
	}

	if w.turn-left.SpawnTurn() < w.settings.PeepSpawnInterval {
		return fmt.Errorf("Too few turns since last spawn for %v", left.ID())
	}

	if w.turn-right.SpawnTurn() < w.settings.PeepSpawnInterval {
		return fmt.Errorf("Too few turns since last spawn for %v", right.ID())
	}

	var locLeft, locRight Location
	var err error
	if locLeft, err = w.ExisterLocation(left); err != nil {
		return fmt.Errorf("Exister %v does not exist...", left)
	}
	if locRight, err = w.ExisterLocation(right); err != nil {
		return fmt.Errorf("Exister %v does not exist...", right)
	}

	newLocation, err := w.FindEmptyLocation(locLeft, locRight)
	if err != nil {
		return fmt.Errorf("Unable to find empty location next to spawners!")
	}

	if rand.Float64() < w.settings.SpawnProbability {
		w.NewPeep(left.Gender(), newLocation)
		left.SetSpawnTurn(w.turn)
		right.SetSpawnTurn(w.turn)
	}
	return nil
}

// DiffGenderSpawn makes a new peep next to one of the provided peeps, if they are of a different gender
func (w *World) DiffGenderSpawn(left, right Exister) error {
	if left.Gender() != right.Gender() {
		var locLeft, locRight Location
		var err error
		if locLeft, err = w.ExisterLocation(left); err != nil {
			return fmt.Errorf("Exister %v does not exist...", left)
		}
		if locRight, err = w.ExisterLocation(right); err != nil {
			return fmt.Errorf("Exister %v does not exist...", right)
		}
		newLocation, err := w.FindEmptyLocation(locLeft, locRight)
		if err == nil {
			if rand.Float64() < w.settings.SpawnProbability {
				w.NewPeep("", newLocation)
				left.SetSpawnTurn(w.turn)
				right.SetSpawnTurn(w.turn)
			}
		}
	}
	return nil
}

// Meet is called when two Existers bump into each other
func (w *World) Meet(left, right Exister) {

	// If they are of the same gender, they spawn a new one (yes yes, I know it's backwards)
	// Spawns only happen the first time peeps meet
	if !left.MetPeep(right) && !right.MetPeep(left) { // no need to check both?
		if err := w.SameGenderSpawn(left, right); err != nil {
			//Log(left.Age(), right.Age())
			//Log(err)
		}
	}
	// Record the meeting
	left.Meet(right, w.turn)
	right.Meet(left, w.turn)

	// If they are of a different gender, they spawn a random child.
	//w.DiffGenderSpawn(left, right)

}

// LocationExister return an exister at the location
func (w *World) LocationExister(l Location) Exister {
	return w.grid.objects.GetByLocation(l)
}

// IsOccupiedLocation returns True if the given Location is occupied by something alive
func (w *World) IsOccupiedLocation(l Location) bool {
	e := w.LocationExister(l)
	if e == nil {
		return false
	}
	if e.IsAlive() {
		return true
	}
	return false
}

// UpdateGrid updates a location on the world grid with the current occupant
// If cell is already occupied, call Meet function and return an error.
// A dead peep is not an occupant
func (w *World) UpdateGrid(e Exister, src Location, dst Location) error {
	if src.SameAs(dst) {
		// Set explicitly again to catch new peeps being created
		w.grid.objects.Set(e, src)
		return nil
	}
	// Check if someone else is already squatting here
	if w.IsOccupiedLocation(dst) {
		squatter := w.grid.objects.GetByLocation(dst)
		if squatter != nil && squatter.ID() != e.ID() {
			// We have a meeting, perhaps something happens here
			w.Meet(e, squatter)
			return fmt.Errorf("Location (%v) taken by: %v", src, squatter.ID())
		}
	}

	// Check that new location is inside the grid
	if w.IsOutsideGrid(dst.X, dst.Y, dst.Z) {
		return fmt.Errorf("Location %v is outside the grid", Location{dst.X, dst.Y, dst.Z})
	}

	w.grid.objects.DelByLocation(src)
	w.grid.objects.Set(e, dst)

	return nil
}

func (w *World) totalNeighbors(l Location, viewDistance int32) int32 {

	// view of 0 means can't see at all
	if viewDistance == 0 {
		return 0
	}

	var neighbors int32

	for x := -viewDistance; x <= viewDistance; x++ {
		for y := -viewDistance; y <= viewDistance; y++ {
			newLoc := NewLocationXYZ(l.X+x, l.Y+y, l.Z)
			if newLoc.SameAs(l) {
				continue // skip our own location
			}
			if !w.IsOutsideGrid(newLoc.X, newLoc.Y, newLoc.Z) {
				neighbors++
			}
		}
	}

	return neighbors
}

// SpawnPoint returns a spawn point for the given Exister
func (w *World) SpawnPoint(e Exister) Location {
	return e.Homebase()
}

// CheckOutsideGrid return error if the move would place object outside grid.
func (w *World) CheckMovementOutsideGrid(src Location, x, y, z int32) error {
	newX := src.X + x
	newY := src.Y + y
	newZ := src.Z + z
	if w.IsOutsideGrid(newX, newY, newZ) {
		return fmt.Errorf("Location %v is outside the grid!", Location{newX, newY, newZ})
	}
	return nil
}

// CheckOutsideGrid returns error if coordinates are outside the grid
// X and Y also remove 1 line for border
func (w *World) IsOutsideGrid(x, y, z int32) bool {
	if x > w.settings.Size.MaxX-1 || x < w.settings.Size.MinX+1 {
		return true
	}
	if y > w.settings.Size.MaxY-1 || y < w.settings.Size.MinY+1 {
		return true
	}
	if z > w.settings.Size.MaxZ || z < w.settings.Size.MinZ {
		return true
	}

	return false
}

// ExisterLocation returns a location, given an exister.
func (w *World) ExisterLocation(e Exister) (Location, error) {
	return w.grid.objects.GetByExister(e)
}

// Move moves a mover in direction and magnitude specified.
// e.g. (1,0,0) will move X to the right 1 and nothing on y and z
func (w *World) Move(e Exister, x, y, z int32) error {
	var src, dst Location
	var err error
	if src, err = w.ExisterLocation(e); err != nil {
		return fmt.Errorf("Exister %v not found on grid.", e)
	}
	if err := w.CheckMovementOutsideGrid(src, x, y, z); err != nil {
		return err
	}

	dst = NewLocationXYZ(src.X+x, src.Y+y, src.Z+z)
	if err := w.UpdateGrid(e, src, dst); err != nil {
		return err
	}
	return nil
}

// ExisterIcon returns the correct icon to use based on criteria defined
func (w *World) ExisterIcon(e Exister) rune {
	midAge := w.settings.SpawnAge

	// icon is the first character of gender
	icon := rune(e.Gender()[0])

	// UpperCase for those who reach middle age
	if e.Age() < midAge {
		return unicode.ToLower(icon)
	}
	return unicode.ToUpper(icon)
}

func colorToTermbox(c PeepGender) termbox.Attribute {
	switch c {
	case "blue":
		return termbox.ColorBlue
	case "red":
		return termbox.ColorRed
	case "green":
		return termbox.ColorGreen
	case "yellow":
		return termbox.ColorYellow
	}
	return termbox.ColorDefault
}

// ExisterFg returns the correct foreground color for an Exister
func (w *World) ExisterFg(e Exister) termbox.Attribute {
	switch e.Gender() {
	case "blue":
		return termbox.ColorBlue
	case "red":
		return termbox.ColorRed
	case "green":
		return termbox.ColorGreen
	case "yellow":
		return termbox.ColorYellow
	}
	return termbox.ColorDefault
}

// ExisterBg returns the correct background color for an Exister
func (w *World) ExisterBg(e Exister) termbox.Attribute {
	// Young ones are highlighted in white < 10 years
	if e.Age() < w.settings.YoungHightlightAge {
		return termbox.ColorWhite
	}

	return termbox.ColorDefault
}

// Visuals describe visual attributes for displaying an Exister
type Visuals struct {
	Char rune              // character displayed
	Fg   termbox.Attribute // foreground color
	Bg   termbox.Attribute // background color
}

// ExisterVisuals returns all the visuals for a given Exister
func (w *World) ExisterVisuals(e Exister) *Visuals {
	v := &Visuals{
		Char: '*',
		Fg:   termbox.ColorDefault,
		Bg:   termbox.ColorDefault,
	}

	v.Char = w.ExisterIcon(e)
	v.Fg = w.ExisterFg(e)
	v.Bg = w.ExisterBg(e)

	return v
}

func (w *World) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w.DrawGrid()

	flashVisuals := &Visuals{
		Char: '☠',
		Fg:   termbox.ColorMagenta,
		Bg:   termbox.ColorBlack,
	}

	for _, loc := range w.grid.objects.AllNonEmptyLocations() {
		e := w.grid.objects.GetByLocation(loc)
		// Convert our coordinates to termbox
		termX := int(loc.X) + int(math.Abs(float64(w.settings.Size.MinX)))
		termY := int(loc.Y) + int(math.Abs(float64(w.settings.Size.MinY)))
		flashForXTurns := Turn(3)

		if !e.IsAlive() {
			// Flash empty squares where peep died for 3 turns
			if w.turn-e.DeadAtTurn() <= flashForXTurns {
				termbox.SetCell(termX, termY, flashVisuals.Char, flashVisuals.Fg, flashVisuals.Bg)
			}
			continue
		}

		visuals := w.ExisterVisuals(e)
		termbox.SetCell(termX, termY, visuals.Char, visuals.Fg, visuals.Bg)
	}
	termbox.Flush()
}

// DrawGrid draws borders around the world and spawn points
func (w *World) DrawGrid() {
	width, height := termbox.Size()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Origin
	termbox.SetCell(0, 0, ' ', termbox.ColorYellow, termbox.ColorYellow)

	// top line
	for x := 0; x <= width-2; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// bottom line
	for x := 0; x <= width-2; x++ {
		termbox.SetCell(x, height-3, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// left border
	for y := 0; y <= height-3; y++ {
		termbox.SetCell(0, y, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// right border
	for y := 0; y <= height-3; y++ {
		termbox.SetCell(width-2, y, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// Homebases
	for gender, loc := range w.homebase {
		termX := int(loc.X) + int(math.Abs(float64(w.settings.Size.MinX)))
		termY := int(loc.Y) + int(math.Abs(float64(w.settings.Size.MinY)))
		termbox.SetCell(termX, termY, ' ', colorToTermbox(gender), colorToTermbox(gender))
	}

	termbox.Flush()
}
